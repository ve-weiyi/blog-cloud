package gormlogx

import (
	"context"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/sdk/trace"
)

// https://zhuanlan.zhihu.com/p/652570748
func TestGoZero(t *testing.T) {
	c := logx.LogConf{
		ServiceName:         "12121",
		Mode:                "file",
		Encoding:            "json",
		TimeFormat:          time.RFC3339,
		Path:                "./logs",
		Level:               "debug",
		MaxContentLength:    0,
		Compress:            false,
		Stat:                false,
		KeepDays:            10,
		StackCooldownMillis: 0,
		MaxBackups:          0,
		MaxSize:             0,
		Rotation:            "",
	}
	logx.SetUp(c)

	ctx := logx.ContextWithFields(context.Background(), logx.Field("trace", "test"))

	handlerName := "gormlogx"
	tracer := trace.NewTracerProvider().Tracer(handlerName)
	ctx, span := tracer.Start(ctx, handlerName)
	defer span.End()

	logx.Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Debug("hello world")
	logx.WithContext(ctx).Slow("hello world")
	logx.WithContext(ctx).Error("hello world")
	logx.Close()
}
