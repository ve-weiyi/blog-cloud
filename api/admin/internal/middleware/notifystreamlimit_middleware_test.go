package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/sse"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
)

// 达到上限时必须拒绝，且返回信封错误——这只能在中间件里做：
// SSE handler 一旦开始写流就已接管响应头，之后发不出信封。
func TestStreamLimitRejectsWhenFull(t *testing.T) {
	broker := sse.NewBroker[notifyx.Event]()
	defer broker.Close()

	m := NewNotifyStreamLimitMiddleware(broker, 2)

	// 占满两个连接（订阅者不读，仅占位）
	_, cancelA := broker.Subscribe()
	defer cancelA()
	_, cancelB := broker.Subscribe()
	defer cancelB()

	reached := false
	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	handler(rr, httptest.NewRequest(http.MethodGet, "/notify-messages/stream", http.NoBody))

	if reached {
		t.Fatal("已达上限时不应放行到 handler")
	}
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("状态码应为 503, got %d", rr.Code)
	}

	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是信封: %v (body=%s)", err, rr.Body.String())
	}
	if body.Code != bizcode.CodeServiceUnavailable {
		t.Fatalf("业务标识应为 %s, got %s", bizcode.CodeServiceUnavailable, body.Code)
	}
}

// 未达上限时放行
func TestStreamLimitAllowsWhenBelowCapacity(t *testing.T) {
	broker := sse.NewBroker[notifyx.Event]()
	defer broker.Close()

	m := NewNotifyStreamLimitMiddleware(broker, 2)

	_, cancel := broker.Subscribe()
	defer cancel()

	reached := false
	handler := m.Handle(func(w http.ResponseWriter, r *http.Request) {
		reached = true
	})

	handler(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/notify-messages/stream", http.NoBody))

	if !reached {
		t.Fatal("未达上限时应放行")
	}
}

// 上限非法时回落到默认值，不能变成"全拒"
func TestStreamLimitFallsBackToDefault(t *testing.T) {
	broker := sse.NewBroker[notifyx.Event]()
	defer broker.Close()

	m := NewNotifyStreamLimitMiddleware(broker, 0)
	if m.limit != DefaultStreamLimit {
		t.Fatalf("非法上限应回落到默认值 %d, got %d", DefaultStreamLimit, m.limit)
	}
}
