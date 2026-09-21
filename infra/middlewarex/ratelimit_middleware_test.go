package middlewarex

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

type fakeLimiter struct {
	result limitx.Result
	err    error
	gotKey string
}

func (f *fakeLimiter) Take(ctx context.Context, key string) (limitx.Result, error) {
	f.gotKey = key

	return f.result, f.err
}

func runRateLimit(t *testing.T, limiter limitx.Limiter, req *http.Request) (*httptest.ResponseRecorder, bool) {
	t.Helper()

	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		responsex.Response(r, w, map[string]any{"ok": true}, nil)
	}

	w := httptest.NewRecorder()
	NewRateLimitMiddleware(limiter).Handle(next)(w, req)

	return w, called
}

func newRequest() *http.Request {
	r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
	r.RemoteAddr = "10.0.0.1:5555"
	r.Header.Set("X-Real-IP", "1.2.3.4")

	return r
}

// TestRateLimitHeadersOnAllowed 守住余量头的给出：客户端要能据此提前降速，
// 而不是撞到 429 才知道。
func TestRateLimitHeadersOnAllowed(t *testing.T) {
	w, called := runRateLimit(t, &fakeLimiter{
		result: limitx.Result{State: limitx.Allowed, Limit: 5, Remaining: 3, Reset: 42 * time.Second},
	}, newRequest())

	if !called {
		t.Fatal("未超限时请求应继续往下走")
	}
	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get(responsex.HeaderRateLimitLimit); got != "5" {
		t.Errorf("RateLimit-Limit = %q, want %q", got, "5")
	}
	if got := w.Header().Get(responsex.HeaderRateLimitRemaining); got != "3" {
		t.Errorf("RateLimit-Remaining = %q, want %q", got, "3")
	}
	if got := w.Header().Get(responsex.HeaderRateLimitReset); got != "42" {
		t.Errorf("RateLimit-Reset = %q, want %q", got, "42")
	}
	if got := w.Header().Get(responsex.HeaderRetryAfter); got != "" {
		t.Errorf("未超限时不应带 Retry-After，得到 %q", got)
	}
}

// TestRateLimitOverQuota 守住拒绝路径：429 + 限流标识 + 精确 Retry-After。
func TestRateLimitOverQuota(t *testing.T) {
	w, called := runRateLimit(t, &fakeLimiter{
		result: limitx.Result{State: limitx.OverQuota, Limit: 5, Remaining: 0, Reset: 12 * time.Second},
	}, newRequest())

	if called {
		t.Fatal("超限时不应继续往下走")
	}
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusTooManyRequests)
	}

	var body struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应体不是 JSON: %v", err)
	}
	if body.Code != bizcode.CodeRateLimited {
		t.Errorf("code = %q, want %q", body.Code, bizcode.CodeRateLimited)
	}

	if got := w.Header().Get(responsex.HeaderRetryAfter); got != "12" {
		t.Errorf("Retry-After = %q, want %q（取限流器给出的精确值）", got, "12")
	}
	if got := w.Header().Get(responsex.HeaderRateLimitRemaining); got != "0" {
		t.Errorf("RateLimit-Remaining = %q, want %q", got, "0")
	}
}

// TestRateLimitFailsOpen 守住降级：限流器报错时放行——限流是保护性约束而非正确性约束，
// 不能因为计数不可得就把全站请求挡在 429 后面。
func TestRateLimitFailsOpen(t *testing.T) {
	w, called := runRateLimit(t, &fakeLimiter{err: errors.New("redis down")}, newRequest())

	if !called {
		t.Fatal("限流器报错时应放行")
	}
	if w.Code != http.StatusOK {
		t.Errorf("HTTP 状态 = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get(responsex.HeaderRateLimitLimit); got != "" {
		t.Errorf("计数不可得时不应给余量头，得到 %q", got)
	}
}

// TestRateLimitKeyIsClientIPAndPath 守住限流维度：按客户端真实 IP + 路径计量，
// 而不是按直连来源（反代地址会退化成一个桶）。
func TestRateLimitKeyIsClientIPAndPath(t *testing.T) {
	limiter := &fakeLimiter{result: limitx.Result{State: limitx.Allowed, Limit: 5, Remaining: 4, Reset: time.Minute}}

	runRateLimit(t, limiter, newRequest())

	if limiter.gotKey != "1.2.3.4:/api/v1/article" {
		t.Errorf("限流 key = %q, want %q", limiter.gotKey, "1.2.3.4:/api/v1/article")
	}
}
