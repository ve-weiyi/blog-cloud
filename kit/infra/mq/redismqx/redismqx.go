package redismqx

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/logz"
)

//有一个分布式系统业务场景：运行有相同代码的三台服务器A、B、C，用户调用接口时可能分配到三台服务器中的任何一台。我希望无论是那一台服务器接收到，都会通知到所有的服务器

// RedisMqConn 结构体，同时实现发布和订阅功能
type RedisMqConn struct {
	logz.Logger

	client  *redis.Client
	channel string
}

// NewRedisMqConn 创建 RedisMqConn 实例
func NewRedisMqConn(redisAddr, channel string) *RedisMqConn {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return &RedisMqConn{
		Logger:  logz.S(),
		client:  client,
		channel: channel,
	}
}

// PublishMessage 发布消息
func (r *RedisMqConn) PublishMessage(ctx context.Context, msg []byte) error {
	err := r.client.Publish(ctx, r.channel, msg).Err()
	if err != nil {
		r.Infof("[发布失败] %v", err)
		return err
	}
	r.Infof("[发布成功] %s\n", msg)
	return nil
}

// SubscribeMessage 订阅消息
func (r *RedisMqConn) SubscribeMessage(handler func(ctx context.Context, msg []byte) error) {
	ctx := context.Background()
	pubsub := r.client.Subscribe(ctx, r.channel)
	defer pubsub.Close()

	r.Infof("[订阅中...]")
	ch := pubsub.Channel()
	for msg := range ch {
		err := handler(ctx, []byte(msg.Payload))
		if err != nil {
			r.Infof("[消息处理失败] %v", err)
		}
	}
}
