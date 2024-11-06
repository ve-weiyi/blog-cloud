package tracex

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/otel/trace"
)

// 如果ctx已经是TraceContext，则返回。否则返回一个新的TraceContext
func MustTraceContext(ctx context.Context) context.Context {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return ctx
	}

	return NewRandomTraceContext()
}

// 创建一个TraceContext
func NewRandomTraceContext() context.Context {
	traceID, err := trace.TraceIDFromHex(generateRandomString(32)) // "4bf92f3577b34da6a3ce929d0e0e4736"
	if err != nil {
		return context.Background()
	}
	spanID, err := trace.SpanIDFromHex(generateRandomString(16)) // "00f067aa0ba902b7"
	if err != nil {
		return context.Background()
	}
	state, err := trace.ParseTraceState("") // "key1=value1,key2=value2"
	if err != nil {
		return context.Background()
	}
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceState: state,
		Remote:     true,
	})
	sc.TraceID()
	ctx := context.Background()
	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	return ctx
}

// 生成指定长度的随机字符串
func generateRandomString(length int) string {
	const charset = "0123456789abcdef"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
