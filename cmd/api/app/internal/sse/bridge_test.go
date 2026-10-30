package sse

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

// 客户端断开后即便没有新信号也必须退出：只阻塞在消息通道上会泄漏 goroutine 与订阅。
func TestForward_StopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan *redis.Message)
	client := make(chan *types.NotificationStreamEvent, 1)
	cancel()

	done := make(chan error, 1)
	go func() { done <- forward(ctx, events, client) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("forward: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("ctx 取消后 forward 未退出")
	}
}

// 订阅收到信号时转成一帧「收件箱变了」，正文由客户端回查。
func TestForward_DeliversSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan *redis.Message, 1)
	client := make(chan *types.NotificationStreamEvent, 1)
	go func() { _ = forward(ctx, events, client) }()

	events <- &redis.Message{}

	select {
	case ev := <-client:
		if ev.Event != notifyx.EventInboxChanged {
			t.Errorf("事件 = %q，期望 %q", ev.Event, notifyx.EventInboxChanged)
		}
	case <-time.After(time.Second):
		t.Fatal("未转发信号")
	}
}
