package rabbitmqx

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/logz"
)

type RabbitmqConn struct {
	url    string
	config *amqp.Config

	connection *amqp.Connection
	channel    *amqp.Channel

	reconnectDelay time.Duration // 重连间隔
	reconnectOnce  sync.Once
	mu             sync.RWMutex
	close          chan struct{}
	closed         bool

	logz.Logger
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
		close:          make(chan struct{}),
		closed:         false,
		Logger:         logz.S(),
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

	r.channel, err = r.connection.Channel()
	if err != nil {
		r.connection.Close()
		return err
	}

	r.reconnectOnce.Do(func() {
		go r.reconnect()
	})

	return nil
}

func (r *RabbitmqConn) reconnect() {
	// recover panic
	defer func() {
		if err := recover(); err != nil {
			r.Infof("Panic: rabbitmq tryReconnect: %v", err)
			debug.PrintStack()
		}
	}()

	for {
		connCloseNotify := r.connection.NotifyClose(make(chan *amqp.Error))
		chanCloseNotify := r.channel.NotifyClose(make(chan *amqp.Error))

		select {
		case err := <-connCloseNotify:
			if err != nil {
				r.Infof("RabbitmqConn connection closed: %v. Attempting to tryReconnect...", err)
				if r.tryReconnect() {
					continue // 重连成功，重新注册 notify
				}
			}
		case err := <-chanCloseNotify:
			if err != nil {
				r.Infof("RabbitmqConn channel closed: %v. Attempting to tryReconnect...", err)
				if r.tryReconnect() {
					continue // 重连成功，重新注册 notify
				}
			}
		case <-r.close:
			return
		}
	}
}

func (r *RabbitmqConn) tryReconnect() bool {
	for {
		err := r.connect()
		if err != nil {
			r.Infof("Failed to tryReconnect. Retrying in %v seconds...", r.reconnectDelay.Seconds())
			time.Sleep(r.reconnectDelay)
			continue
		}
		r.Info("Reconnected to RabbitmqConn")
		return true
	}
}

func (r *RabbitmqConn) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return
	}
	r.closed = true

	if r.channel != nil {
		r.channel.Close()
	}

	if r.connection != nil {
		r.connection.Close()
	}

	close(r.close)
	r.Infof("RabbitmqConn connection closed")
}

func (r *RabbitmqConn) Declare(q *QueueOptions, e *ExchangeOptions, b *BindingOptions) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := declareQueue(r.channel, q); err != nil {
		return fmt.Errorf("declare queue failed: %w", err)
	}
	if err := declareExchange(r.channel, e); err != nil {
		return fmt.Errorf("declare exchange failed: %w", err)
	}
	if err := declareBindings(r.channel, q, e, b); err != nil {
		return fmt.Errorf("declare bindings failed: %w", err)
	}
	return nil
}

func (r *RabbitmqConn) Consume(opts ConsumerOptions) (<-chan amqp.Delivery, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
		tableToAMQPTable(opts.Args),
	)
}

func (r *RabbitmqConn) Publish(msg amqp.Publishing, opts PublisherOptions) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

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
