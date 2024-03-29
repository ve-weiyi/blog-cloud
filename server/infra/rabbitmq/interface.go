package rabbitmq

type MessagePublisher interface {
	// 发送消息
	PublishMessage(message []byte) error
}

type MessageSubscriber interface {
	// 接收消息
	SubscribeMessage(handler func(message []byte) error) error
}
