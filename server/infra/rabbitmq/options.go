package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// ExchangeOptions rabbitmq 交换机
type ExchangeOptions struct {
	// 交换机明
	Name string
	// 交换机类型
	Type ExchangeType
	// 是否持久化
	Durable bool
}

type QueueOptions struct {
	// 消息队列名。如果为空rabbitmq服务器回创建一个唯一的队列名
	Name string
	// 是否持久化 rabbitmq服务重启后，队列数据不会丢失；消费者连接时，队列也不会被删除
	Durable bool
	// 参数
	Args amqp.Table
}

type Options struct {
	// 是否自动确认消息
	AutoAck bool

	// 设置当autoAck为false时，消费失败将消息重入队列
	Requeue bool

	// 路由key
	// 在pub/sub模式下，这里的key要为空
	// 其中“*”用于匹配一个单词，“#”用于匹配多个单词（可以是零个）
	Key string

	// 队列配置
	Queue QueueOptions

	// 交换机配置
	Exchange ExchangeOptions
	// 是否需要申明交换机，true表示不需要
	WithoutExchange bool
}

type Option func(options *Options)

func newOptions(opts ...Option) Options {
	// 初始订阅相关配置。默认开启自动ack
	opt := Options{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// DisableAutoAck 禁止自动确认消息
func DisableAutoAck() Option {
	return func(o *Options) {
		o.AutoAck = false
	}
}

// Requeue 设置消费失败时将消息重入队列
func Requeue() Option {
	return func(o *Options) {
		o.Requeue = true
	}
}

// Key 设置路由key
func Key(key string) func(o *Options) {
	return func(o *Options) {
		o.Key = key
	}
}

// Queue 设置队列名
func Queue(q QueueOptions) Option {
	return func(o *Options) {
		o.Queue = q
	}
}

// Exchange 设置交换机配置
func Exchange(ex ExchangeOptions) Option {
	return func(o *Options) {
		o.Exchange = ex
	}
}

// WithoutExchange 不申明交换机
func WithoutExchange() Option {
	return func(o *Options) {
		o.WithoutExchange = true
	}
}
