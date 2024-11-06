package gormlogx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/logger"
)

// New initialize logger.
func New(config logger.Config) logger.Interface {
	return &traceLogger{
		Config: config,
	}
}

type traceLogger struct {
	logger.Config
}

// LogMode log mode.
func (l *traceLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return l
}

// Info print info.
func (l *traceLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logx.WithContext(ctx).WithCallerSkip(4).Infow(fmt.Sprintf(msg, data...))
	}
}

// Warn print warn messages.
func (l *traceLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logx.WithContext(ctx).WithCallerSkip(4).Errorw(fmt.Sprintf(msg, data...))
	}
}

// Error print error messages.
func (l *traceLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logx.WithContext(ctx).WithCallerSkip(4).Errorw(fmt.Sprintf(msg, data...))
	}
}

// Trace print sql message.
func (l *traceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		logx.WithContext(ctx).WithCallerSkip(4).Errorw("gorm error",
			logx.Field("rows", rows),
			logx.Field("duration", timeDuration(elapsed)),
			logx.Field("sql", sql),
			logx.Field("error", fmt.Sprintf("%+v", err)))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		logx.WithContext(ctx).WithCallerSkip(4).Errorw("gorm slow info",
			logx.Field("rows", rows),
			logx.Field("duration", timeDuration(elapsed)),
			logx.Field("sql", sql))
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		logx.WithContext(ctx).WithCallerSkip(4).Infow("gorm info",
			logx.Field("rows", rows),
			logx.Field("duration", timeDuration(elapsed)),
			logx.Field("sql", sql))
	}
}

func timeDuration(duration time.Duration) string {
	return fmt.Sprintf("%.3fms", float32(duration)/float32(time.Millisecond))
}
