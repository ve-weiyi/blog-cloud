// Package mq 是消息队列的传输层与事件契约。
//
// 本包只依赖 config 与 vkit；消费者在 mqlogic 子包（它依赖 svc），
// 本包不得反向依赖 svc —— 否则与 svc → mq 成环。
package mq

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/infra/metrics"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/adapter/mqx/event"
	"github.com/ve-weiyi/vkit/adapter/mqx/impl/rabbitmqx"
)

// Client 持有事件运行时与共享发布器。
//
// 启动策略：只有 MQ 允许降级 —— RabbitMQ 不可用时记录错误并继续启动，
// 发布返回 ErrUnavailable、订阅等待重连；其余依赖（DB/Redis/邮件/短信）失败即失败。
type Client struct {
	Runtime   *event.Runtime
	Publisher mqx.Publisher
}

// NewClient 建立 RabbitMQ 连接与事件运行时。
func NewClient(c config.RabbitMQConf) *Client {
	logger := newSlogLogger()
	broker, err := rabbitmqx.NewBroker(&rabbitmqx.Config{
		URL:          rabbitmqx.BuildAMQPURL(c.Host, c.Port, c.Username, c.Password),
		ExchangeType: rabbitmqx.ExchangeTopic,
		LazyConnect:  true,
		Logger:       logger,
		Metrics:      metrics.NewMQX(),
	})
	if err != nil {
		logx.Errorf("[mq] create broker: %v", err)
		return &Client{}
	}

	runtime := event.NewRuntime(broker, logger, event.WithMessagingSystem("rabbitmq"))
	publisher, err := runtime.NewPublisher()
	if err != nil {
		logx.Errorf("[mq] create publisher: %v", err)
	}
	return &Client{Runtime: runtime, Publisher: publisher}
}

// Close 先停止订阅、等待在途处理，再关闭连接。
func (c *Client) Close() {
	if c == nil || c.Runtime == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Runtime.Close(ctx); err != nil {
		logx.Errorf("[mq] close: %v", err)
	}
}
