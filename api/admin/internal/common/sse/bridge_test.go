package sse

import (
	"context"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
)

// fakeSource 只提供一条事件通道，用于不引入 Broker 的场景
type fakeSource struct {
	ch chan notifyx.Event
}

func (f *fakeSource) Subscribe() (<-chan notifyx.Event, func()) {
	return f.ch, func() {}
}

func TestToStreamEventMapsAllFields(t *testing.T) {
	in := notifyx.Event{
		Type:        notifyx.EventNotice,
		MessageId:   5,
		Title:       "维护通知",
		Category:    "maintenance",
		Level:       "warning",
		PublishedAt: 1758000000,
	}

	got := ToStreamEvent(in)

	want := types.NotifyStreamEvent{
		Event:       notifyx.EventNotice,
		MessageId:   5,
		Title:       "维护通知",
		Category:    "maintenance",
		Level:       "warning",
		PublishedAt: 1758000000,
	}
	if *got != want {
		t.Fatalf("映射不符:\n got=%+v\nwant=%+v", *got, want)
	}
}

// 事件必须原样送到推送流上
func TestStreamForwardsEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := &fakeSource{ch: make(chan notifyx.Event, 1)}
	client := make(chan *types.NotifyStreamEvent, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client)
	}()

	src.ch <- notifyx.Event{Type: notifyx.EventNoticeRevoke, MessageId: 7}

	select {
	case frame := <-client:
		if frame.Event != notifyx.EventNoticeRevoke || frame.MessageId != 7 {
			t.Fatalf("推送帧不符: %+v", frame)
		}
	case <-time.After(time.Second):
		t.Fatal("未收到推送帧")
	}

	cancel()
	if err := <-done; err != nil {
		t.Fatalf("取消后应正常返回: %v", err)
	}
}

// 真实广播器的事件也能桥接过来（覆盖 Subscribe 的实际行为）
func TestStreamForwardsBrokerEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := NewBroker[notifyx.Event]()
	defer broker.Close()
	client := make(chan *types.NotifyStreamEvent, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, broker, client)
	}()

	// 订阅建立之前发布的事件会丢（那时广播器还没有订阅者），
	// 所以周期性重发直到收到帧
	var frame *types.NotifyStreamEvent
	deadline := time.After(2 * time.Second)
	for frame == nil {
		broker.Publish(notifyx.Event{Type: notifyx.EventNotice, MessageId: 1})
		select {
		case f := <-client:
			frame = f
		case <-time.After(20 * time.Millisecond):
		case <-deadline:
			t.Fatal("未收到推送帧")
		}
	}

	if frame.Event != notifyx.EventNotice || frame.MessageId != 1 {
		t.Fatalf("推送帧不符: %+v", frame)
	}

	cancel()
	<-done
}

// 消费端不再读取（客户端已断）时，桥接必须靠 ctx 退出，不能永久阻塞在写通道上
func TestStreamExitsWhenConsumerStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	src := &fakeSource{ch: make(chan notifyx.Event, 1)}
	client := make(chan *types.NotifyStreamEvent) // 无缓冲且无人读取

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client)
	}()

	// 让桥接协程阻塞在写通道上
	src.ch <- notifyx.Event{Type: notifyx.EventNotice}

	select {
	case err := <-done:
		t.Fatalf("尚在运行时不应返回: %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("应当正常返回: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("ctx 取消后仍未退出——写通道阻塞导致 goroutine 泄漏")
	}
}

// 事件源关闭时桥接结束
func TestStreamEndsWhenSourceCloses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := &fakeSource{ch: make(chan notifyx.Event)}
	client := make(chan *types.NotifyStreamEvent, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client)
	}()

	close(src.ch)

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("事件源关闭后应正常返回: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("事件源关闭后桥接未退出")
	}
}
