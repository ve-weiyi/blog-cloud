package sse

import (
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
)

func TestHandleDeliversDecodedEventToBroker(t *testing.T) {
	broker := NewBroker[notifyx.Event]()
	defer broker.Close()
	sub := NewSubscriber(nil, broker)

	ch, cancel := broker.Subscribe()
	defer cancel()

	payload, err := notifyx.Event{
		Type:      notifyx.EventNotice,
		MessageId: 9,
		Title:     "系统维护",
	}.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	sub.handle(payload)

	got, _ := recv(t, ch)
	if got.MessageId != 9 || got.Type != notifyx.EventNotice || got.Title != "系统维护" {
		t.Fatalf("投递的事件不符: %+v", got)
	}
}

// 载荷格式错误只应被丢弃，不得 panic 或中断订阅
func TestHandleIgnoresInvalidPayload(t *testing.T) {
	broker := NewBroker[notifyx.Event]()
	defer broker.Close()
	sub := NewSubscriber(nil, broker)

	ch, cancel := broker.Subscribe()
	defer cancel()

	sub.handle("not-json")

	select {
	case e := <-ch:
		t.Fatalf("非法载荷不应投递: %+v", e)
	case <-time.After(50 * time.Millisecond):
	}
}
