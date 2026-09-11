package mq

import (
	"fmt"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/adapter/mqx/rabbitmqx"
)

const (
	LogoutExchange   = "blog-logout-exchange"
	LogoutQueue      = "blog-logout-queue"
	LogoutRoutingKey = "logout"
)

func initLogoutMq(svcCtx *svc.ServiceContext) (mqx.MessageQueue, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		svcCtx.Config.RabbitMQConf.Username,
		svcCtx.Config.RabbitMQConf.Password,
		svcCtx.Config.RabbitMQConf.Host,
		svcCtx.Config.RabbitMQConf.Port)

	return rabbitmqx.NewRabbitMQ(&rabbitmqx.Config{
		URL:          url,
		ExchangeName: LogoutExchange,
		ExchangeType: "fanout",
		Durable:      true,
		AutoDelete:   false,
		QueueConfig: &rabbitmqx.QueueConfig{
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
		},
	})
}
