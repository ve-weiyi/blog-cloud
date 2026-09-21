package sse

// 订阅必须是**进程级一份**：Redis 只有一条订阅流，各 SSE 连接从广播器分食。

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
)

// Subscriber 订阅通知事件 channel，把事件投进本实例的通知广播器。
type Subscriber struct {
	rdb    *redis.Client
	broker *Broker[notifyx.Event]
}

// NewSubscriber 创建订阅器
func NewSubscriber(rdb *redis.Client, broker *Broker[notifyx.Event]) *Subscriber {
	return &Subscriber{rdb: rdb, broker: broker}
}

// Start 在后台保持订阅，订阅断开后隔 3 秒重试，直到 ctx 结束。
func (s *Subscriber) Start(ctx context.Context) {
	if s.rdb == nil {
		logx.Error("通知事件订阅器缺少 Redis 客户端，推送功能不可用")
		return
	}

	go func() {
		for {
			if err := s.subscribe(ctx); err != nil && !errors.Is(err, context.Canceled) {
				logx.Errorf("通知事件订阅中断: %v, 3 秒后重连", err)
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
			}
		}
	}()
}

func (s *Subscriber) subscribe(ctx context.Context) error {
	sub := s.rdb.Subscribe(ctx, cachekey.NotifyEventChannel)
	defer func() {
		_ = sub.Close()
	}()

	ch := sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return context.Canceled
		case msg, ok := <-ch:
			if !ok {
				return errors.New("通知事件订阅通道已关闭")
			}
			if msg.Channel != cachekey.NotifyEventChannel {
				continue
			}
			s.handle(msg.Payload)
		}
	}
}

// handle 解析载荷并投入广播器。
//
// 载荷非法时静默跳过——事件是尽力而为的提醒，不该因一条坏消息中断订阅。
// 订阅者缓冲写满导致的丢帧由广播器自行处理，客户端重连回查时会补齐。
func (s *Subscriber) handle(payload string) {
	event, err := notifyx.Decode(payload)
	if err != nil {
		logx.Errorf("解析通知事件失败（载荷 %d 字节）: %v", len(payload), err)
		return
	}
	s.broker.Publish(event)
}
