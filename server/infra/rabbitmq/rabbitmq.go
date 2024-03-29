package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrRabbitmqClosed = fmt.Errorf("rabbitmq closed")

var (
	defaultHeartbeat  = 5 * time.Second
	defaultAmqpConfig = amqp.Config{
		Heartbeat: defaultHeartbeat,
	}
)

// ExchangeType 交换机类型
type ExchangeType string

const (
	ExchangeTypeSimple = ""       //Simple模式（简单模式，一个生产者对应多个消费者）  只需要声明队列
	ExchangeTypeWorker = ""       //Worker模式（工作模式，一个生产者对应多个消费者，一个消息只能被一个消费者获取） 只需要声明队列
	ExchangeTypeFanout = "fanout" //Publish模式（订阅模式，消息被路由投递给多个队列，一个消息被多个消费者获取） 只需要声明交换机
	ExchangeTypeDirect = "direct" //Direct模式（路由模式，消息被路由投递给符合路由规则的队列，一个消息被一个消费者获取） 声明交换机和key
	ExchangeTypeTopic  = "topic"  //Topic模式（主题模式，消息被路由投递给符合通配符的队列，一个消息被一个消费者获取） 声明交换机和key
)

type RabbitmqConn struct {
	url        string
	connection *amqp.Connection
	channel    *amqp.Channel

	opts Options
	// 内部状态维护 section
	sync.Mutex
	connected       bool
	close           chan struct{}
	waitConnection  chan struct{}
	connCloseNotify chan *amqp.Error
	chanCloseNotify chan *amqp.Error
}

func NewRabbitmqConn(url string, opts ...Option) *RabbitmqConn {
	ret := &RabbitmqConn{
		url:            url,
		waitConnection: make(chan struct{}),
		close:          make(chan struct{}),
		opts:           newOptions(opts...),
	}
	close(ret.waitConnection)

	return ret
}

func (r *RabbitmqConn) Connect(config *amqp.Config) error {
	r.Lock()
	if r.connected {
		r.Unlock()
		return nil
	}

	r.Unlock()

	if config == nil {
		config = &defaultAmqpConfig
	}
	return r.connect(config)
}

func (r *RabbitmqConn) connect(config *amqp.Config) error {
	if err := r.tryConnect(config); err != nil {
		return err
	}

	r.Lock()
	r.connected = true
	r.Unlock()

	// 创建重连协程
	go r.reconnect(config)

	return nil
}

func (r *RabbitmqConn) tryConnect(config *amqp.Config) (err error) {
	if r.connection == nil || r.connection.IsClosed() {
		r.connection, err = amqp.DialConfig(r.url, *config)
		if err != nil {
			return
		}

		r.connCloseNotify = make(chan *amqp.Error, 1)
		r.connection.NotifyClose(r.connCloseNotify)

		r.channel, err = r.connection.Channel()
		if err != nil {
			return
		}

		r.chanCloseNotify = make(chan *amqp.Error, 1)
		r.channel.NotifyClose(r.chanCloseNotify)
	}

	err = r.DeclareQueue(r.opts.Queue)
	if err != nil {
		return
	}

	// 不需要交换机
	if !r.opts.WithoutExchange {
		// 创建交换机
		err = r.DeclareExchange(r.opts.Exchange)
		if err != nil {
			return
		}

		// 队列与交换机绑定
		err = r.BindQueue(r.opts.Queue.Name, r.opts.Key, r.opts.Exchange.Name, nil)
		if err != nil {
			return err
		}
	}

	log.Printf("rabbitmq 连接成功！运行模式:'%v',交换机:'%v',队列:'%v',key:'%v'", r.opts.Exchange.Type, r.opts.Exchange.Name, r.opts.Queue.Name, r.opts.Key)
	return
}

func (r *RabbitmqConn) reconnect(config *amqp.Config) {
	// recover panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic: rabbitmq reconnect: %v", err)
			debug.PrintStack()
		}
	}()

	for {
		if !r.connected { // 第一次 connected 为 true, 不需要重连
			b := NewForeverBackoff()
			for {
				err := r.tryConnect(config)
				if err == nil {
					break
				}

				// 等待重试
				log.Printf("Error: rabbitmq reconnect: %v", err)
				b.Wait()
				continue
			}

			r.Lock()
			r.connected = true
			r.Unlock()

			// 通知等待的协程链接创建好了
			close(r.waitConnection)
			log.Println("Info: rabbitmq reconnect success")
		}

		// 监听关闭事件
		select {
		case <-r.close:
			return

		case err := <-r.connCloseNotify: // 连接关闭通知
			log.Printf("Warning: connection notify close: %v", err)
			r.Lock()
			r.connected = false
			r.waitConnection = make(chan struct{})
			r.Unlock()
			r.connCloseNotify = nil

		case err := <-r.chanCloseNotify: // channel关闭通知
			log.Printf("Warning: channel notify close: %v", err)
			r.Lock()
			r.connected = false
			r.waitConnection = make(chan struct{})
			r.Unlock()
			r.chanCloseNotify = nil
		}
	}
}

func (r *RabbitmqConn) Close() error {
	r.Lock()
	defer r.Unlock()

	select {
	case <-r.close: // 已经关闭了直接返回
		return nil
	default:
		close(r.close)
		r.connected = false
	}

	r.channel.Close()
	return r.connection.Close()
}

func (r *RabbitmqConn) DeclareExchange(ex ExchangeOptions) error {
	return r.channel.ExchangeDeclare(
		ex.Name,         // name
		string(ex.Type), // kind
		ex.Durable,      // durable
		false,           // autoDelete
		//true表示这个exchange不可以被client用来推送消息，仅用来进行exchange和exchange之间的绑定
		false, // internal
		false, // noWait
		nil,   // args
	)
}

// DeclareQueue 非持久化队列
// rabbitmq服务重启后，队列数据会丢失；消费者连接时，队列会自动删除
func (r *RabbitmqConn) DeclareQueue(qx QueueOptions) (err error) {
	_, err = r.channel.QueueDeclare(
		qx.Name,    // name
		qx.Durable, // durable
		//是否自动删除
		false, // autoDelete
		//是否具有排他性
		false, // internal
		//是否阻塞处理
		false,   // noWait
		qx.Args, // args
	)

	return
}

func (r *RabbitmqConn) BindQueue(queue, key, exchange string, args amqp.Table) error {
	return r.channel.QueueBind(
		queue,    // name
		key,      // key
		exchange, // exchange
		false,    // noWait
		args,     // args
	)
}

// Publish 发布消息，如果发送失败会自动重试
// 注意：发消息前一定确保交换机成功创建。交换机不存在 amqp 并不会返回错误，而是直接关闭 rabbitmq channel
func (r *RabbitmqConn) PublishMessage(message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-r.close:
		return ErrRabbitmqClosed
	case <-ctx.Done():
		return ctx.Err()
	case <-r.waitConnection:
	}

	if r.channel == nil {
		return fmt.Errorf("channel is nil")
	}

	var key string
	switch r.opts.Exchange.Type {
	case ExchangeTypeFanout, ExchangeTypeDirect, ExchangeTypeTopic:
		key = r.opts.Key
	default:
		key = r.opts.Queue.Name
	}

	return r.channel.PublishWithContext(
		ctx,
		r.opts.Exchange.Name,
		key,
		//如果为true，根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

// Subscribe 订阅消息
func (r *RabbitmqConn) SubscribeMessage(handler func(message []byte) error) error {
	if r.channel == nil {
		return fmt.Errorf("channel is nil")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	// 绑定消费者
	deliveries, err := r.channel.Consume(
		r.opts.Queue.Name, // queue
		id.String(),       // consumer
		r.opts.AutoAck,    // autoAck
		false,             // exclusive
		false,             // noLocal
		false,             // nowait
		nil,               // args
	)
	if err != nil {
		log.Printf("Error: rabbitmq subscribe consume: %v", err)
		return err
	}

	fn := func(msg amqp.Delivery) {
		err = handler(msg.Body)
		if err != nil {
			log.Printf("Waring: rabbitmq subscribe handle: %v", err)
		}

		// 自动ack直接返回
		if r.opts.AutoAck {
			return
		}

		// 消费成功ack
		if err == nil {
			if err = msg.Ack(false); err != nil {
				log.Printf("Warning: rabbitmq ack msg: %v", err)
			}
			return
		}

		// 消费失败，根据配置决定是否重入队列
		if err = msg.Nack(false, r.opts.Requeue); err != nil {
			log.Printf("Warning: rabbitmq nack msg: %v", err)
		}
	}

	forever := make(chan bool)
	go func() {
		// recover panic
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: rabbitmq subscribe: %v", err)
				debug.PrintStack()
			}
		}()

		// 无限循环执行
		for {
			select {
			case <-r.close: // 我们主动关闭，退出程序
				return
			case <-time.After(time.Second): // 收到断连通知后 r.conn.waitConnection 会被重新赋值，这里是为了防止出现死锁
				continue
			case <-r.waitConnection: // 如果在正在重连，等待重连成功
			}

			r.Lock()
			if !r.connected {
				r.Unlock()
				continue
			}

			r.Unlock()

			// 当 rabbitmq conn 和 channel 关闭时 deliveries 都会被关闭，
			// 所以这里会退出循环，等待重连后重新订阅。
			// 除以上两个情况外，就是队列被删除了，deliveries 也会被关闭，但是并不会收任何关闭通知消息。
			for d := range deliveries {
				fn(d)
			}
		}
	}()

	<-forever
	return nil
}
