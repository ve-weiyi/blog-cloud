package outbox

import (
	"context"
	"testing"
)

// 调用方已有的头优先；返回新 map，不修改入参。
func TestWithTraceHeader_CopiesAndPreserves(t *testing.T) {
	headers := map[string]string{"trace_id": "existing", "k": "v"}

	out := withTraceHeader(context.Background(), headers)
	if out["trace_id"] != "existing" {
		t.Fatalf("trace_id = %q, want existing", out["trace_id"])
	}
	if out["k"] != "v" {
		t.Fatalf("k = %q, want v", out["k"])
	}

	out["mutated"] = "yes"
	if _, ok := headers["mutated"]; ok {
		t.Fatal("不应修改调用方传入的 headers")
	}
}

// 没有 trace 时也要返回可序列化的空对象，而不是 nil。
func TestWithTraceHeader_NoTrace(t *testing.T) {
	out := withTraceHeader(context.Background(), nil)
	if out == nil {
		t.Fatal("应返回非 nil map")
	}
	if len(out) != 0 {
		t.Fatalf("无 trace 时不应写入头，实际 %v", out)
	}
}
