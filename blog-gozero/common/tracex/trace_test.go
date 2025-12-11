package tracex

import (
	"context"
	"testing"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestContext(t *testing.T) {
	// 往context添加一个key
	ctx := context.WithValue(context.Background(), "trace_id", "value001")
	// 打印日志
	logx.WithContext(ctx).Info("我的日志1,没有trace_id")

	// context注入traceID
	tracer := trace.NewTracerProvider().Tracer("测试节点")
	ctx, span := tracer.Start(ctx, "test")
	defer span.End()

	// 打印日志
	logx.WithContext(ctx).Info("我的日志2,有trace_id")

}
