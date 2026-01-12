package logz

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

func TestZap(t *testing.T) {
	L().Info("hello world")
}

func TestTrace(t *testing.T) {
	ip := "localhost"
	ctx := context.Background()
	tracer := trace.NewTracerProvider().Tracer(ip)
	ctx, span := tracer.Start(ctx, TraceIDKey)
	defer span.End()

	SetLog(&LogConfig{
		Level:      "info",
		Filename:   "./test.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
	})
	L().Info("Trace test", zap.String("ip", ip))
	L().Info("Trace test", WithTraceField(ctx))
}
