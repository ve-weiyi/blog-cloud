package outbox

import (
	"context"

	"github.com/zeromicro/go-zero/core/trace"
)

const traceHeader = "trace_id"

// withTraceHeader 把当前 trace 写入消息头，返回新的 map。
//
// outbox 是延迟发布：relay 发布时的 ctx 已不是原始请求的 ctx，不在这里
// 捕获的话，异步段就与请求链路断开了。调用方已有的值优先。
func withTraceHeader(ctx context.Context, headers map[string]string) map[string]string {
	out := make(map[string]string, len(headers)+1)
	for key, value := range headers {
		out[key] = value
	}
	if out[traceHeader] != "" {
		return out
	}
	if id := trace.TraceIDFromContext(ctx); id != "" {
		out[traceHeader] = id
	}
	return out
}
