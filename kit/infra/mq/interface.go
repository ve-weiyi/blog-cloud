package mq

import (
	"context"
)

type MessagePublisher interface {
	// 发送消息
	PublishMessage(ctx context.Context, msg []byte) error
}

type MessageSubscriber interface {
	// 接收消息
	SubscribeMessage(handler func(ctx context.Context, msg []byte) error)
}
