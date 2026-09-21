package limitx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// fakeEvalClient 只回放预设的 Lua 返回值，用于验证 Go 侧对返回值的解码。
type fakeEvalClient struct {
	reply interface{}
	err   error
}

func (f *fakeEvalClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd {
	return redis.NewCmdResult(f.reply, f.err)
}

// TestPeriodLimitDecodesResult 守住 Lua 返回值到 Result 的解码：
// 状态、余量、重置时长三者都要如实映射，不能错位。
func TestPeriodLimitDecodesResult(t *testing.T) {
	cases := []struct {
		name          string
		reply         interface{}
		wantState     int
		wantRemaining int
		wantReset     time.Duration
	}{
		{
			name:          "配额未用尽",
			reply:         []interface{}{int64(Allowed), int64(3), int64(60_000)},
			wantState:     Allowed,
			wantRemaining: 3,
			wantReset:     time.Minute,
		},
		{
			name:          "刚好用尽",
			reply:         []interface{}{int64(HitQuota), int64(0), int64(30_000)},
			wantState:     HitQuota,
			wantRemaining: 0,
			wantReset:     30 * time.Second,
		},
		{
			name:          "超出配额",
			reply:         []interface{}{int64(OverQuota), int64(0), int64(12_000)},
			wantState:     OverQuota,
			wantRemaining: 0,
			wantReset:     12 * time.Second,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewPeriodLimit(60, 5, nil, "test:")
			l.client = &fakeEvalClient{reply: c.reply}

			got, err := l.Take(context.Background(), "key")
			if err != nil {
				t.Fatalf("Take 返回错误: %v", err)
			}
			if got.State != c.wantState {
				t.Errorf("State = %d, want %d", got.State, c.wantState)
			}
			if got.Limit != 5 {
				t.Errorf("Limit = %d, want 5", got.Limit)
			}
			if got.Remaining != c.wantRemaining {
				t.Errorf("Remaining = %d, want %d", got.Remaining, c.wantRemaining)
			}
			if got.Reset != c.wantReset {
				t.Errorf("Reset = %v, want %v", got.Reset, c.wantReset)
			}
		})
	}
}

// TestPeriodLimitPropagatesError 守住错误上抛：Redis 失败要把错误交给调用方，
// 由调用方决定取舍（限流中间件的取舍是放行并告警）。
func TestPeriodLimitPropagatesError(t *testing.T) {
	l := NewPeriodLimit(60, 5, nil, "test:")
	l.client = &fakeEvalClient{err: errors.New("redis down")}

	if _, err := l.Take(context.Background(), "key"); err == nil {
		t.Fatal("Redis 失败时 Take 应返回错误")
	}
}

// TestPeriodLimitRejectsMalformedReply 守住返回值校验：Lua 返回值形状不对时报错，
// 不臆造成"放行"或"拒绝"。
func TestPeriodLimitRejectsMalformedReply(t *testing.T) {
	cases := []struct {
		name  string
		reply interface{}
	}{
		{"返回值不是数组", int64(1)},
		{"返回值元素个数不对", []interface{}{int64(1), int64(2)}},
		{"返回值元素不是整数", []interface{}{"x", "y", "z"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := NewPeriodLimit(60, 5, nil, "test:")
			l.client = &fakeEvalClient{reply: c.reply}

			if _, err := l.Take(context.Background(), "key"); err == nil {
				t.Fatal("返回值形状不对时 Take 应返回错误")
			}
		})
	}
}
