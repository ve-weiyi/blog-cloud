package ssex

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
)

// testFrame 是测试用的推送帧，与任何服务的 types 无关。
type testFrame struct {
	Type string
	ID   int64
}

func toTestFrame(event notifyx.Event) testFrame {
	return testFrame{Type: event.Type, ID: event.AnnouncementId}
}

// fakeSource 只提供一条事件通道，用于不引入 Broker 的场景
type fakeSource struct {
	ch chan notifyx.Event
}

func (f *fakeSource) Subscribe() (<-chan notifyx.Event, func()) {
	return f.ch, func() {}
}

// 事件必须经 mapEvent 转换后原样送到推送流上
func TestBridgeForwardsMappedEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := &fakeSource{ch: make(chan notifyx.Event, 1)}
	client := make(chan testFrame, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client, toTestFrame)
	}()

	src.ch <- notifyx.Event{Type: notifyx.EventAnnouncementRevoked, AnnouncementId: 7}

	select {
	case frame := <-client:
		if frame.Type != notifyx.EventAnnouncementRevoked || frame.ID != 7 {
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
func TestBridgeForwardsBrokerEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := NewBroker[notifyx.Event]()
	defer broker.Close()
	client := make(chan testFrame, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, broker, client, toTestFrame)
	}()

	// 订阅建立之前发布的事件会丢（那时广播器还没有订阅者），
	// 所以周期性重发直到收到帧
	var frame testFrame
	deadline := time.After(2 * time.Second)
	for frame == (testFrame{}) {
		broker.Publish(notifyx.Event{Type: notifyx.EventAnnouncementPublished, AnnouncementId: 1})
		select {
		case frame = <-client:
		case <-time.After(20 * time.Millisecond):
		case <-deadline:
			t.Fatal("未收到推送帧")
		}
	}

	if frame.Type != notifyx.EventAnnouncementPublished || frame.ID != 1 {
		t.Fatalf("推送帧不符: %+v", frame)
	}

	cancel()
	<-done
}

// 消费端不再读取（客户端已断）时，桥接必须靠 ctx 退出，不能永久阻塞在写通道上
func TestBridgeExitsWhenConsumerStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	src := &fakeSource{ch: make(chan notifyx.Event, 1)}
	client := make(chan testFrame) // 无缓冲且无人读取

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client, toTestFrame)
	}()

	// 让桥接协程阻塞在写通道上
	src.ch <- notifyx.Event{Type: notifyx.EventAnnouncementPublished}

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
func TestBridgeEndsWhenSourceCloses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	src := &fakeSource{ch: make(chan notifyx.Event)}
	client := make(chan testFrame, 4)

	done := make(chan error, 1)
	go func() {
		done <- Bridge(ctx, src, client, toTestFrame)
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

// 客户端断开后即便没有新信号也必须退出：只阻塞在消息通道上会泄漏 goroutine 与订阅。
func TestForwardStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan *redis.Message)
	client := make(chan testFrame, 1)
	cancel()

	done := make(chan error, 1)
	go func() { done <- forward(ctx, events, client, func() testFrame { return testFrame{} }) }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("forward: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("ctx 取消后 forward 未退出")
	}
}

// 每收到一条信号写一帧，帧内容由调用方决定，正文由客户端回查。
func TestForwardDeliversSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan *redis.Message, 1)
	client := make(chan testFrame, 1)
	go func() { _ = forward(ctx, events, client, func() testFrame { return testFrame{Type: "inbox"} }) }()

	events <- &redis.Message{}

	select {
	case frame := <-client:
		if frame.Type != "inbox" {
			t.Errorf("帧 = %+v，期望 Type=inbox", frame)
		}
	case <-time.After(time.Second):
		t.Fatal("未转发信号")
	}
}
