package job_test

import (
	"context"
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/job"
)

// 单轮 panic 不应终止循环，后续轮次照常执行。
func TestPeriodicRunsAfterPanic(t *testing.T) {
	runs := make(chan struct{}, 4)
	periodic := job.NewPeriodic("test", 5*time.Millisecond, time.Millisecond, func(context.Context) {
		select {
		case runs <- struct{}{}:
		default:
		}
		panic("boom")
	})
	periodic.Start()

	for round := 1; round <= 2; round++ {
		select {
		case <-runs:
		case <-time.After(2 * time.Second):
			t.Fatalf("第 %d 轮未执行", round)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	periodic.Stop(ctx)
	periodic.Stop(ctx) // 重复停止无副作用
}
