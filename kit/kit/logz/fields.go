package logz

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const TraceIDKey = "traceId"

// 从 Context 中提取 TraceID 并生成 Zap Field
func WithTraceField(ctx context.Context) zap.Field {
	spanCtx := trace.SpanContextFromContext(ctx)
	return zap.String(TraceIDKey, spanCtx.TraceID().String())
}

func WitheContext(ctx context.Context) *zap.Logger {
	// 从 Context 中提取 TraceID 并生成 Zap Field
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		return defaultLogger.With(zap.String(TraceIDKey, spanCtx.TraceID().String()))
	}

	return defaultLogger
}
