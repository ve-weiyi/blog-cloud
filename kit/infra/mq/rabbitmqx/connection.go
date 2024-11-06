package rabbitmqx

import (
	"log"
	"runtime/debug"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitmqConn struct {
	url    string
	config *amqp.Config

	connection *amqp.Connection
	channel    *rabbitmqChannel

	reconnectDelay time.Duration // 重连间隔
	close          chan bool
	mu             sync.Mutex
}

func NewRabbitmqConn(url string, config *amqp.Config) (*RabbitmqConn, error) {
	if config == nil {
		config = &amqp.Config{
			Heartbeat: 5 * time.Second,
		}
	}

	r := &RabbitmqConn{
		url:    url,
		config: config,

		reconnectDelay: 5 * time.Second,
		close:          make(chan bool),
	}

	if err := r.connect(); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RabbitmqConn) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var err error
	r.connection, err = amqp.DialConfig(r.url, *r.config)
	if err != nil {
		return err
	}

	r.channel, err = newRabbitChannel(r.connection)
	if err != nil {
		return err
	}

	go r.reconnect()

	return nil
}

func (r *RabbitmqConn) reconnect() {
	// recover panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic: rabbitmq tryReconnect: %v", err)
			debug.PrintStack()
		}
	}()

	connCloseNotify := r.connection.NotifyClose(make(chan *amqp.Error))
	chanCloseNotify := r.channel.rawChannel.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case err := <-connCloseNotify:
			if err != nil {
				log.Printf("RabbitmqConn connection closed: %v. Attempting to tryReconnect...", err)
				r.tryReconnect()
			}
		case err := <-chanCloseNotify:
			if err != nil {
				log.Printf("RabbitmqConn channel closed: %v. Attempting to tryReconnect...", err)
				r.tryReconnect()
			}
		case <-r.close: // 监听到关闭信号，退出重连循环
			return
		}
	}
}

func (r *RabbitmqConn) tryReconnect() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for {
		err := r.connect()
		if err != nil {
			log.Printf("Failed to tryReconnect. Retrying in %v seconds...", r.reconnectDelay.Seconds())
			time.Sleep(r.reconnectDelay)
		}
		log.Println("Reconnected to RabbitmqConn")
		return
	}
}

func (r *RabbitmqConn) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		r.channel.Close()
	}

	if r.connection != nil {
		r.connection.Close()
	}

	r.close <- true
	log.Println("RabbitmqConn connection closed")
}

func (r *RabbitmqConn) Declare(q *QueueOptions, e *ExchangeOptions, b *BindingOptions) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.channel.Declare(q, e, b)
}

func (r *RabbitmqConn) Consume(opts ConsumerOptions) (<-chan amqp.Delivery, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return nil, amqp.ErrClosed
	}

	return r.channel.Consume(
		opts.Queue,
		opts.Name,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
}

func (r *RabbitmqConn) Publish(msg amqp.Publishing, opts PublisherOptions) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel == nil {
		return amqp.ErrClosed
	}

	return r.channel.Publish(
		opts.Exchange,
		opts.RoutingKey,
		opts.Mandatory,
		opts.Immediate,
		msg)
}
