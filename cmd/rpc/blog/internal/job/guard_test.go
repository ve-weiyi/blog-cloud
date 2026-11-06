package job_test

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/job"
)

// 单轮 panic 不得终止调用方，后续工作照常执行。
func TestGuardRecovers(t *testing.T) {
	job.Guard("test", func() { panic("boom") })

	ran := false
	job.Guard("test", func() { ran = true })
	if !ran {
		t.Fatal("Guard 之后的调用应继续执行")
	}
}
