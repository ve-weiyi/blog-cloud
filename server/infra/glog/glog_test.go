package glog

import (
	"context"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/glog/zaplog"
)

// 同时输出到控制台和文件测试
func TestGlogConsoleAndFile(t *testing.T) {
	cfg := zaplog.NewFileConfig()
	cfg.Encoding = zaplog.FormatConsole
	cfg.Path = "./logs"
	cfg.FileName = "blog.log"
	cfg.Prefix = "W "
	cfg.Level = "Info"
	cfg.ShowColor = true

	logger := NewGlogger(1, cfg)

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
}

// 压力测试
func BenchmarkExample(b *testing.B) {
	logs := zaplog.NewZapLogger(0, zaplog.NewConsoleConfig())
	logj := zaplog.NewZapLogger(0, zaplog.NewFileConfig())
	for i := 0; i < b.N; i++ {
		// 在这里执行需要进行压力测试的代码
		logs.Sugar().Info("hello", i)
		logj.Sugar().Info("hello", i)
	}
}

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

	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Info("hello world")
	logx.WithContext(ctx).Debug("hello world")
	logx.WithContext(ctx).Slow("hello world")
	logx.WithContext(ctx).Error("hello world")
	logx.Close()
}
