package notifyx

import (
	"encoding/json"
	"testing"
)

// 事件载荷是 rpc 与 api 网关之间的跨进程契约，字段名一旦漂移，
// 接收侧就会解出空事件——所以这里同时锁住 JSON 键名与往返一致性。
func TestEventEncodeDecodeRoundTrip(t *testing.T) {
	in := Event{
		Type:        EventNotice,
		MessageId:   42,
		Title:       "系统维护",
		Category:    "maintenance",
		Level:       "warning",
		PublishedAt: 1758000000,
	}

	payload, err := in.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	out, err := Decode(payload)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	if out != in {
		t.Fatalf("往返不一致:\n in=%+v\nout=%+v", in, out)
	}
}

func TestEventEncodeUsesSnakeCaseKeys(t *testing.T) {
	payload, err := Event{Type: EventNoticeRevoke, MessageId: 7}.Encode()
	if err != nil {
		t.Fatalf("Encode: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if raw["event"] != EventNoticeRevoke {
		t.Fatalf("event 字段缺失或值不符: %v", raw)
	}
	if _, ok := raw["message_id"]; !ok {
		t.Fatalf("message_id 键名漂移: %v", raw)
	}
}

func TestDecodeInvalidPayloadReturnsError(t *testing.T) {
	if _, err := Decode("not-json"); err == nil {
		t.Fatal("非法载荷应当返回错误")
	}
}
