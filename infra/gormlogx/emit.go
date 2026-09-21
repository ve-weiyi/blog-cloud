package gormlogx

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
	"gorm.io/gorm/logger"
)

// kind 是 Trace 对一次数据库操作的归类，决定它记不记、记成哪一类。
type kind int

const (
	kindSkip      kind = iota // 不记录
	kindQuery                 // 普通查询
	kindMigration             // 结构变更语句（DDL）
	kindSlow                  // 超过慢阈值
	kindFailure               // 执行出错

	kindCount // 哨兵，不是有效的 kind，只用来给 emission 定长度
)

// logLevel 是本包往 go-zero 发日志时可用的级别。
//
// logx 只有 debug / info / error / slow 四档，没有 warn —— 这是 gorm 的
// Warn 被映射到 error 的原因。
type logLevel int

const (
	levelInfo logLevel = iota
	levelError
	levelSlow
)

// emission 是 kind 到发射方式的映射。改级别或改文案只动这张表。
//
// 用数组而不是 map，是为了让索引在编译期就有界：新增 kind 却忘记登记时，
// 拿到的是一个零值表项，而不是一次查表失败。
var emission = [kindCount]struct {
	level logLevel
	msg   string
}{
	kindSkip:      {}, // 不会被发射，Trace 在此之前已经返回
	kindQuery:     {level: levelInfo, msg: "gorm query"},
	kindMigration: {level: levelSlow, msg: "gorm migration"},
	kindSlow:      {level: levelError, msg: "gorm slow"},
	kindFailure:   {level: levelError, msg: "gorm error"},
}

// 日志字段名。统一加 gorm. 前缀，避免与 go-zero 的内建字段撞名 ——
// logx 自己就用 duration 这个 key 输出时长，直接叫 duration 会出现两个同名字段。
const (
	fieldSQL      = "gorm.sql"
	fieldRows     = "gorm.rows"
	fieldDuration = "gorm.duration"
	fieldError    = "gorm.error"
)

// callerSkip 是交给 go-zero WithCallerSkip 的值，用于让日志的 caller 字段
// 指向业务代码而不是本包内部。
//
// go-zero 内部取调用者的表达式是 getCaller(callerDepth + callerSkip)，其中
// callerDepth 固定为 4。重构前本包直接从各方法里发日志，该值取 4；抽出 emit
// 之后调用链多了一帧，故取 5。
//
// 两点提醒：
//   - 在本包内再增删层数都会平移这个值，而且不会报错，只会让 caller 静默指错位置。
//   - 它无法用单元测试标定：测试的调用栈太浅，go-zero 的 getCaller 会返回空串。
//     核对方式是看真实运行的日志。
const callerSkip = 5

// classify 判定一次数据库操作该不该记录、归到哪一类。
//
// 纯函数：不读取 SQL、不依赖 logx，可独立做表驱动测试。
// 优先级照搬 gorm 原生 logger：出错 > 超过慢阈值 > 记录。
func classify(cfg logger.Config, elapsed time.Duration, err error) kind {
	switch {
	case err != nil && cfg.LogLevel >= logger.Error &&
		(!errors.Is(err, logger.ErrRecordNotFound) || !cfg.IgnoreRecordNotFoundError):
		return kindFailure
	case elapsed > cfg.SlowThreshold && cfg.SlowThreshold != 0 && cfg.LogLevel >= logger.Warn:
		return kindSlow
	case cfg.LogLevel == logger.Info:
		return kindQuery
	default:
		return kindSkip
	}
}

// isDDL 判断一条 SQL 是否是结构变更语句。gorm 发出的 SQL 均为大写，故只认大写前缀。
func isDDL(sql string) bool {
	s := strings.TrimLeft(sql, " \t\n\r")

	return strings.HasPrefix(s, "ALTER") ||
		strings.HasPrefix(s, "CREATE") ||
		strings.HasPrefix(s, "DROP")
}

// queryFields 组装一次查询的日志字段。
func queryFields(sql string, rows int64, elapsed time.Duration, err error) []logx.LogField {
	fields := []logx.LogField{
		logx.Field(fieldSQL, sql),
		logx.Field(fieldRows, rows),
		// 复用 go-zero 自己的时长表示，与本仓其它日志保持同一格式
		logx.Field(fieldDuration, timex.ReprOfDuration(elapsed)),
	}
	if err != nil {
		fields = append(fields, logx.Field(fieldError, err.Error()))
	}

	return fields
}

// emit 是全包唯一往 logx 发日志的出口。
func emit(ctx context.Context, level logLevel, msg string, fields ...logx.LogField) {
	l := logx.WithContext(ctx).WithCallerSkip(callerSkip)

	switch level {
	case levelError:
		l.Errorw(msg, fields...)
	case levelSlow:
		l.Sloww(msg, fields...)
	default:
		// 未登记的级别按 info 处理：对日志来说，级别记低好过整条丢失
		l.Infow(msg, fields...)
	}
}
