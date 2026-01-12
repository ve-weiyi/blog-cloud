package rabbitmqx

// ExchangeType 交换机类型
type ExchangeType string

const (
	//Publish模式（发布/订阅模式，消息被路由投递给多个队列，一个消息被多个消费者获取）
	//发布者：交换机
	//订阅者：交换机、队列
	ExchangeTypeFanout = "fanout"

	//Direct模式（路由模式，消息被路由投递给符合路由规则的队列，一个消息被一个消费者获取）
	//发布者：交换机、key
	//订阅者：交换机、队列、key
	ExchangeTypeDirect = "direct"

	//Topic模式（主题模式，消息被路由投递给符合通配符的队列，一个消息被一个消费者获取）
	//发布者：交换机、key
	//订阅者：交换机、队列、key
	ExchangeTypeTopic = "topic"
)

type QueueOptions struct {
	Name       string //不填是随机队列名称
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       Table
}

type ExchangeOptions struct {
	Name       string
	Kind       string // possible values: empty string for default exchange or direct, topic, fanout
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Passive    bool // if false, a missing exchange will be created on the server
	Args       Table
}

type BindingOptions struct {
	RoutingKey string // fanout模式不需要routing key
	NoWait     bool
	Args       Table
}

type PublisherOptions struct {
	Exchange   string
	RoutingKey string

	Mandatory bool // 确保消息成功路由到一个队列，否则返回给生产者。
	Immediate bool // 确保消息立即投递给消费者，如果队列中没有消费者，消息将被退回给生产者。 （已废弃）
}

type PublisherOption func(*PublisherOptions)

func WithPublisherExchange(exchange string) PublisherOption {
	return func(o *PublisherOptions) {
		o.Exchange = exchange
	}
}

func WithPublisherRoutingKey(routingKey string) PublisherOption {
	return func(o *PublisherOptions) {
		o.RoutingKey = routingKey
	}
}

func WithPublisherMandatory(mandatory bool) PublisherOption {
	return func(o *PublisherOptions) {
		o.Mandatory = mandatory
	}
}

type ConsumerOptions struct {
	Queue string

	Name      string
	AutoAck   bool
	Exclusive bool
	NoWait    bool
	NoLocal   bool
	Args      Table
}

type ConsumerOption func(*ConsumerOptions)

func WithConsumerQueue(queue string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.Queue = queue
	}
}

func WithConsumerName(name string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.Name = name
	}
}

func WithConsumerAutoAck(autoAck bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.AutoAck = autoAck
	}
}

type Options struct {
	// 交换机
	ExchangeOptions *ExchangeOptions
	// 队列
	QueueOptions *QueueOptions
	// 绑定
	BindingOptions *BindingOptions
}

// Option 是一个函数类型，用来设置配置参数
type Option func(*Options)
