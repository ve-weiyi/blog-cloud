package gormlogx

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

// New 返回一个把日志交给 go-zero logx 的 gorm logger，实现 logger.Interface。
//
// 日志出口是 go-zero 的包级 writer 与级别闸门：若 logx 的级别被设置得高于本包
// 要发的级别，日志会被静默丢弃；logx 处于 disableLevel 时连 SetWriter 都会被
// 忽略且不报错。排查「日志没出来」时先确认这一点。
func New(config logger.Config) logger.Interface {
	return &traceLogger{
		Config: config,
	}
}

// traceLogger 实现 gorm 的 logger.Interface，把 gorm 的日志级别映射到 logx。
type traceLogger struct {
	logger.Config
}

var _ logger.Interface = (*traceLogger)(nil)

// LogMode 返回一个级别被改写的新 logger，原实例不受影响 ——
// 否则并发复用同一个 logger 会互相污染。
func (l *traceLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level

	return &newLogger
}

// Info 记录 gorm 的普通信息。
func (l *traceLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		emit(ctx, levelInfo, fmt.Sprintf(msg, data...))
	}
}

// Warn 记录 gorm 的警告。logx 没有 warn 级别，故落到 error。
func (l *traceLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		emit(ctx, levelError, fmt.Sprintf(msg, data...))
	}
}

// Error 记录 gorm 的错误。
func (l *traceLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		emit(ctx, levelError, fmt.Sprintf(msg, data...))
	}
}

// Trace 记录一次 SQL 执行。
//
// 分两步：先由 classify 判定该不该记（这一步不需要 SQL），确认要记时才调用 fc
// 做 SQL 格式化 —— 后者不便宜，不该为不输出的查询付出。
func (l *traceLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	k := classify(l.Config, elapsed, err)
	if k == kindSkip {
		return
	}

	sql, rows := fc()
	if k == kindQuery && isDDL(sql) {
		k = kindMigration
	}

	e := emission[k]
	emit(ctx, e.level, e.msg, queryFields(sql, rows, elapsed, err)...)
}
