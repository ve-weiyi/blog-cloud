package mqlogic

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx/event"
)

func TestRegisterRejectsDuplicateHandlers(t *testing.T) {
	router := event.NewRouter(nil)
	if err := Register(router, &svc.ServiceContext{}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := Register(router, &svc.ServiceContext{}); err == nil {
		t.Fatal("同一 routing key 重复注册必须失败")
	}
}
