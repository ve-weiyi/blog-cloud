// Package mq 是消息队列的传输层与事件契约。
package mq

import (
	"context"
	"fmt"
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
// 构造成功即保证 Runtime 与 Publisher 非 nil。运行期的 broker 抖动不在这里兜底：
// 发布返回 ErrUnavailable，订阅断开由 Runtime 自行重建。
type Client struct {
	Runtime   *event.Runtime
	Publisher mqx.Publisher
}

// NewClient 建立 RabbitMQ 连接与事件运行时。
func NewClient(c config.RabbitMQConf) (*Client, error) {
	logger := newSlogLogger()
	broker, err := rabbitmqx.NewBroker(&rabbitmqx.Config{
		URL:          rabbitmqx.BuildAMQPURL(c.Host, c.Port, c.Username, c.Password),
		ExchangeType: rabbitmqx.ExchangeTopic,
		LazyConnect:  true,
		Logger:       logger,
		Metrics:      metrics.NewMQX(),
	})
	if err != nil {
		return nil, fmt.Errorf("mq: create broker: %w", err)
	}

	runtime := event.NewRuntime(broker, logger, event.WithMessagingSystem("rabbitmq"))
	publisher, err := runtime.NewPublisher()
	if err != nil {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = runtime.Close(closeCtx)
		return nil, fmt.Errorf("mq: create publisher: %w", err)
	}
	return &Client{Runtime: runtime, Publisher: publisher}, nil
}

// Close 先停止订阅、等待在途处理，再关闭连接。
func (c *Client) Close() {
	if c == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := c.Runtime.Close(ctx); err != nil {
		logx.Errorf("[mq] close: %v", err)
	}
}
