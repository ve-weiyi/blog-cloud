package sse

import (
	"testing"
	"time"
)

// 用自定义类型做测试，确保泛型对非基础类型同样成立
type tick struct {
	Seq  int
	Name string
}

func recv[T any](t *testing.T, ch <-chan T) (T, bool) {
	t.Helper()
	select {
	case v, ok := <-ch:
		return v, ok
	case <-time.After(time.Second):
		t.Fatal("等待事件超时")
		var zero T
		return zero, false
	}
}

func expectNoEvent[T any](t *testing.T, ch <-chan T) {
	t.Helper()
	select {
	case v, ok := <-ch:
		t.Fatalf("不应收到事件: %+v (open=%v)", v, ok)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestPublishFansOutToAllSubscribers(t *testing.T) {
	b := NewBroker[tick]()
	defer b.Close()

	a, cancelA := b.Subscribe()
	defer cancelA()
	c, cancelC := b.Subscribe()
	defer cancelC()

	want := tick{Seq: 1, Name: "alpha"}
	b.Publish(want)

	if got, _ := recv(t, a); got != want {
		t.Fatalf("订阅者 a 收到 %+v, want %+v", got, want)
	}
	if got, _ := recv(t, c); got != want {
		t.Fatalf("订阅者 c 收到 %+v, want %+v", got, want)
	}
}

func TestLenCountsSubscribers(t *testing.T) {
	b := NewBroker[tick]()

	if b.Len() != 0 {
		t.Fatalf("初始订阅者数应为 0, got %d", b.Len())
	}

	_, cancelA := b.Subscribe()
	_, cancelB := b.Subscribe()
	if b.Len() != 2 {
		t.Fatalf("订阅者数应为 2, got %d", b.Len())
	}

	cancelA()
	if b.Len() != 1 {
		t.Fatalf("注销后订阅者数应为 1, got %d", b.Len())
	}

	cancelA() // 幂等
	if b.Len() != 1 {
		t.Fatalf("重复注销不应改变订阅者数, got %d", b.Len())
	}
	cancelB()
}

// 注销后通道被关闭，消费方以 ok=false 退出（与 scan 的 Broker 一致）
func TestCancelClosesSubscriberChannel(t *testing.T) {
	b := NewBroker[tick]()
	defer b.Close()

	ch, cancel := b.Subscribe()
	cancel()

	if _, ok := <-ch; ok {
		t.Fatal("注销后通道应当已关闭")
	}

	// 注销后再发布不得 panic（写锁与读锁互斥）
	b.Publish(tick{Seq: 9})
}

// 慢客户端（缓冲写满）只丢自己的帧，不影响其他订阅者
func TestSlowSubscriberDoesNotBlockOthers(t *testing.T) {
	b := NewBroker[tick]()
	defer b.Close()

	slow, cancelSlow := b.Subscribe()
	defer cancelSlow()
	_ = slow

	// 先塞满慢订阅者的缓冲
	for i := 0; i < SubscriberBuffer; i++ {
		b.Publish(tick{Seq: i})
	}

	// 缓冲满之后再加一个即时排空的订阅者：只该丢慢的那个
	fast, cancelFast := b.Subscribe()
	defer cancelFast()

	for i := 0; i < 3; i++ {
		seq := SubscriberBuffer + i
		b.Publish(tick{Seq: seq})
		if got, _ := recv(t, fast); got.Seq != seq {
			t.Fatalf("快订阅者收到 %d, want %d", got.Seq, seq)
		}
	}
}

func TestCloseClosesAllSubscriberChannels(t *testing.T) {
	b := NewBroker[tick]()
	ch, _ := b.Subscribe()

	b.Close()

	if _, ok := <-ch; ok {
		t.Fatal("Close 后订阅者通道应当已关闭")
	}
	if b.Len() != 0 {
		t.Fatalf("Close 后订阅者数应为 0, got %d", b.Len())
	}
}

func TestPublishAfterCloseIsNoop(t *testing.T) {
	b := NewBroker[tick]()
	b.Close()

	b.Publish(tick{Seq: 1}) // 不得 panic

	ch, cancel := b.Subscribe() // Close 后订阅：立即关闭，消费方随即退出
	defer cancel()
	if _, ok := <-ch; ok {
		t.Fatal("Close 后订阅的通道应当已关闭")
	}
}

// Close 之后再注销订阅：通道已被 Close 关闭，注销不得再次 close（否则 panic）
func TestCancelAfterCloseDoesNotPanic(t *testing.T) {
	b := NewBroker[tick]()
	ch, cancel := b.Subscribe()

	b.Close()

	// 消费方读到通道关闭后退出（bridge 的正常退出路径），随后执行 defer 里的注销
	for range ch {
	}

	cancel()
}
