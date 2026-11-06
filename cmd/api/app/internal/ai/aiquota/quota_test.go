package aiquota

import (
	"context"
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/config"
)

func TestSubject(t *testing.T) {
	if got := Subject(aictx.Identity{UserId: "u1"}); got != "user:u1" {
		t.Fatalf("user subject = %q", got)
	}
	if got := Subject(aictx.Identity{DeviceId: "d1", Guest: true}); got != "guest:d1" {
		t.Fatalf("guest subject = %q", got)
	}
}

func TestLimit(t *testing.T) {
	conf := config.AiQuotaConf{GuestDaily: 20, UserDaily: 100}
	if got := Limit(conf, true); got != 20 {
		t.Fatalf("guest limit = %d", got)
	}
	if got := Limit(conf, false); got != 100 {
		t.Fatalf("user limit = %d", got)
	}
}

func TestKeys(t *testing.T) {
	if got := DailyKey("user:u1", "20260101"); got != "blog:ai:quota:user:u1:20260101" {
		t.Fatalf("daily key = %q", got)
	}
	if got := IPDailyKey("1.2.3.4", "20260101"); got != "blog:ai:ip:1.2.3.4:20260101" {
		t.Fatalf("ip key = %q", got)
	}
}

func TestUsedWithoutRedis(t *testing.T) {
	used, err := Used(context.Background(), nil, "any")
	if err != nil || used != 0 {
		t.Fatalf("nil redis should read zero, got %d / %v", used, err)
	}
}
