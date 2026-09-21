package gormlogx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/logger"
)

// 上游 trace/span，取值与 vkit 原诊断程序 cmd/mysqltest 一致。
const (
	testTraceID = "4bf92f3577b34da6a3ce929d0e0e4736"
	testSpanID  = "00f067aa0ba902b7"
)

// captured 记录一条被 logx 写出的日志。
type captured struct {
	level  string
	msg    string
	fields map[string]any
}

// captureWriter 实现 logx.Writer，按级别把日志收进内存。
//
// 直接实现 Writer 而不是把文本塞进 buffer，是为了能断言「走了哪个级别」——
// Infow / Errorw / Sloww 最终落到同一个 Writer，只看文本无法区分。
type captureWriter struct {
	mu      sync.Mutex
	records []captured
}

func (w *captureWriter) add(level string, v any, fields ...logx.LogField) {
	w.mu.Lock()
	defer w.mu.Unlock()

	m := make(map[string]any, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	w.records = append(w.records, captured{level: level, msg: fmt.Sprint(v), fields: m})
}

// Alert / Severe / Stack 本包不会用到，实现出来只为满足 logx.Writer。
func (w *captureWriter) Alert(any)                       {}
func (w *captureWriter) Severe(any)                      {}
func (w *captureWriter) Stack(any)                       {}
func (w *captureWriter) Close() error                    { return nil }
func (w *captureWriter) Debug(v any, f ...logx.LogField) { w.add("debug", v, f...) }
func (w *captureWriter) Info(v any, f ...logx.LogField)  { w.add("info", v, f...) }
func (w *captureWriter) Error(v any, f ...logx.LogField) { w.add("error", v, f...) }
func (w *captureWriter) Slow(v any, f ...logx.LogField)  { w.add("slow", v, f...) }
func (w *captureWriter) Stat(v any, f ...logx.LogField)  { w.add("stat", v, f...) }

func (w *captureWriter) snapshot() []captured {
	w.mu.Lock()
	defer w.mu.Unlock()

	return append([]captured(nil), w.records...)
}

// redirect 把 go-zero 的日志出口换成内存捕获器，并返回它。
//
// logx 的级别与 Writer 都是进程级全局状态，故本包测试一律不使用 t.Parallel。
func redirect(t *testing.T) *captureWriter {
	t.Helper()

	w := &captureWriter{}
	// 必须先放开级别：logx 处于 disableLevel 时 SetWriter 会被静默忽略。
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(w)

	t.Cleanup(func() {
		logx.SetWriter(logx.NewWriter(os.Stderr))
	})

	return w
}

// testConfig 造一份用例用的 logger.Config，只填用例关心的字段。
func testConfig(level logger.LogLevel, ignoreNotFound bool) logger.Config {
	return logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: ignoreNotFound,
	}
}

// withSpanContext 造一个带远端 span 的 ctx，模拟上游传下来的 trace。
// 与原诊断程序 cmd/mysqltest 的 withTraceContext 同构。
func withSpanContext(t *testing.T) context.Context {
	t.Helper()

	traceID, err := trace.TraceIDFromHex(testTraceID)
	if err != nil {
		t.Fatalf("TraceIDFromHex(%q): %v", testTraceID, err)
	}
	spanID, err := trace.SpanIDFromHex(testSpanID)
	if err != nil {
		t.Fatalf("SpanIDFromHex(%q): %v", testSpanID, err)
	}

	return trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
		Remote:  true,
	}))
}

func TestNewAndLogMode(t *testing.T) {
	original, ok := New(testConfig(logger.Info, false)).(*traceLogger)
	if !ok {
		t.Fatalf("New 返回值类型 = %T, want *traceLogger", New(testConfig(logger.Info, false)))
	}
	if original.LogLevel != logger.Info {
		t.Errorf("New 未带上配置的级别：got %v, want %v", original.LogLevel, logger.Info)
	}

	changed, ok := original.LogMode(logger.Silent).(*traceLogger)
	if !ok {
		t.Fatalf("LogMode 返回值类型 = %T, want *traceLogger", original.LogMode(logger.Silent))
	}
	if changed == original {
		t.Error("LogMode 应返回新实例")
	}
	// 原实例的级别不得被改动，否则并发复用同一个 logger 会互相污染。
	if original.LogLevel != logger.Info {
		t.Errorf("原实例的级别被改动：got %v, want %v", original.LogLevel, logger.Info)
	}
	if changed.LogLevel != logger.Silent {
		t.Errorf("新实例的级别 = %v, want %v", changed.LogLevel, logger.Silent)
	}
}

// TestLogLevelGating 固定 Info / Warn / Error 三个方法的级别闸门。
// 注意 Warn 走的是 Errorw，因此它在日志里表现为 error 级别。
func TestLogLevelGating(t *testing.T) {
	tests := []struct {
		name             string
		level            logger.LogLevel
		info, warn, erro bool
	}{
		{name: "Silent 三个都不出", level: logger.Silent},
		{name: "Error 只出错误", level: logger.Error, erro: true},
		{name: "Warn 出慢查询与错误", level: logger.Warn, warn: true, erro: true},
		{name: "Info 全出", level: logger.Info, info: true, warn: true, erro: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := redirect(t)
			ctx := context.Background()
			l := New(testConfig(tc.level, false))

			l.Info(ctx, "from-info")
			l.Warn(ctx, "from-warn")
			l.Error(ctx, "from-error")

			byMsg := make(map[string]captured, len(w.snapshot()))
			for _, r := range w.snapshot() {
				byMsg[r.msg] = r
			}

			for _, c := range []struct {
				msg  string
				want bool
				lvl  string
			}{
				{msg: "from-info", want: tc.info, lvl: "info"},
				{msg: "from-warn", want: tc.warn, lvl: "error"},
				{msg: "from-error", want: tc.erro, lvl: "error"},
			} {
				got, ok := byMsg[c.msg]
				if ok != c.want {
					t.Errorf("%s 是否输出 = %v, want %v", c.msg, ok, c.want)
					continue
				}
				if ok && got.level != c.lvl {
					t.Errorf("%s 的级别 = %q, want %q", c.msg, got.level, c.lvl)
				}
			}
		})
	}
}

// TestTrace 覆盖 Trace 的四个分支与 Silent 短路。
func TestTrace(t *testing.T) {
	const selectSQL = "SELECT * FROM user_account WHERE id = ?"

	tests := []struct {
		name           string
		level          logger.LogLevel
		ignoreNotFound bool
		elapsed        time.Duration
		err            error
		sql            string
		wantLog        bool
		wantLevel      string
		wantMsg        string
	}{
		{
			name:    "Silent 级别短路，完全不出日志",
			level:   logger.Silent,
			elapsed: time.Millisecond,
			sql:     selectSQL,
		},
		{
			name:      "普通查询走 info",
			level:     logger.Info,
			elapsed:   time.Millisecond, // 小于 SlowThreshold(200ms)
			sql:       selectSQL,
			wantLog:   true,
			wantLevel: "info",
			wantMsg:   "gorm query",
		},
		{
			name:      "DDL 走 slow",
			level:     logger.Info,
			elapsed:   time.Millisecond,
			sql:       "ALTER TABLE user_account ADD COLUMN nickname varchar(64)",
			wantLog:   true,
			wantLevel: "slow",
			wantMsg:   "gorm migration",
		},
		{
			name:      "超出慢阈值走 error",
			level:     logger.Info,
			elapsed:   time.Second, // 大于 SlowThreshold(200ms)
			sql:       selectSQL,
			wantLog:   true,
			wantLevel: "error",
			wantMsg:   "gorm slow",
		},
		{
			name:      "出错走 error",
			level:     logger.Info,
			elapsed:   time.Millisecond,
			err:       errors.New("syntax error"),
			sql:       selectSQL,
			wantLog:   true,
			wantLevel: "error",
			wantMsg:   "gorm error",
		},
		{
			name:           "IgnoreRecordNotFoundError=false 时记录未找到要打出来",
			level:          logger.Info,
			ignoreNotFound: false,
			elapsed:        time.Millisecond,
			err:            logger.ErrRecordNotFound,
			sql:            selectSQL,
			wantLog:        true,
			wantLevel:      "error",
			wantMsg:        "gorm error",
		},
		{
			// 开关只关掉 error 分支：查询本身仍在 info 分支留下轨迹，
			// 被吞掉的是「这次没查到」这件事，不是这次查询。
			name:           "IgnoreRecordNotFoundError=true 时记录未找到不再报错，只剩 info 轨迹",
			level:          logger.Info,
			ignoreNotFound: true,
			elapsed:        time.Millisecond,
			err:            logger.ErrRecordNotFound,
			sql:            selectSQL,
			wantLog:        true,
			wantLevel:      "info",
			wantMsg:        "gorm query",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			w := redirect(t)
			l := New(testConfig(tc.level, tc.ignoreNotFound))

			l.Trace(context.Background(), time.Now().Add(-tc.elapsed), func() (string, int64) {
				return tc.sql, 1
			}, tc.err)

			got := w.snapshot()
			if !tc.wantLog {
				if len(got) != 0 {
					t.Fatalf("不该出日志，却出了 %d 条：%+v", len(got), got)
				}
				return
			}

			if len(got) != 1 {
				t.Fatalf("日志条数 = %d, want 1：%+v", len(got), got)
			}
			if got[0].level != tc.wantLevel {
				t.Errorf("级别 = %q, want %q", got[0].level, tc.wantLevel)
			}
			if got[0].msg != tc.wantMsg {
				t.Errorf("文案 = %q, want %q", got[0].msg, tc.wantMsg)
			}
		})
	}
}

// TestTraceFields 断言 Trace 带出的字段，含 go-zero 从 ctx 注入的 trace / span。
// 断言用的 trace、span 键名取 go-zero logx 的默认值（defaultTraceKey / defaultSpanKey）。
func TestTraceFields(t *testing.T) {
	const (
		querySQL = "SELECT * FROM user_account WHERE id = ?"
		rows     = int64(7)
	)

	w := redirect(t)
	l := New(testConfig(logger.Info, false))
	l.Trace(withSpanContext(t), time.Now().Add(-3*time.Millisecond), func() (string, int64) {
		return querySQL, rows
	}, nil)

	got := w.snapshot()
	if len(got) != 1 {
		t.Fatalf("日志条数 = %d, want 1：%+v", len(got), got)
	}

	rec := got[0]
	if rec.msg != "gorm query" {
		t.Errorf("文案 = %q, want %q", rec.msg, "gorm query")
	}
	// 字段名按字面量断言：它们是日志的对外形态，改名应当被人看见
	if v := rec.fields["gorm.rows"]; v != rows {
		t.Errorf("gorm.rows = %v(%T), want %v(%T)", v, v, rows, rows)
	}
	if v := rec.fields["gorm.sql"]; v != querySQL {
		t.Errorf("gorm.sql = %v, want %v", v, querySQL)
	}
	if v, ok := rec.fields["gorm.duration"].(string); !ok || !strings.HasSuffix(v, "ms") {
		t.Errorf("gorm.duration = %v, want 形如 \"1.2ms\" 的字符串", rec.fields["gorm.duration"])
	}
	// trace / span 由 go-zero 的 logx 从 ctx 注入，键名取 go-zero 的默认值
	if v := rec.fields["trace"]; v != testTraceID {
		t.Errorf("trace = %v, want %v", v, testTraceID)
	}
	if v := rec.fields["span"]; v != testSpanID {
		t.Errorf("span = %v, want %v", v, testSpanID)
	}
}

// kindNames 只给测试用，让断言失败时打印出类别名而不是数字。
var kindNames = map[kind]string{
	kindSkip:      "skip",
	kindQuery:     "query",
	kindMigration: "migration",
	kindSlow:      "slow",
	kindFailure:   "failure",
}

// TestClassify 固定「该不该记、归到哪一类」的判定。
// classify 是纯函数：不碰 logx，因此这里不需要 redirect，也不会与其它用例争全局状态。
func TestClassify(t *testing.T) {
	tests := []struct {
		name            string
		level           logger.LogLevel
		elapsed         time.Duration
		err             error
		ignoreNF        bool
		noSlowThreshold bool
		want            kind
	}{
		{
			name:    "Silent 级别一律不记",
			level:   logger.Silent,
			elapsed: time.Second,
			want:    kindSkip,
		},
		{
			name:    "零值级别降级为不记",
			level:   0,
			elapsed: time.Second,
			want:    kindSkip,
		},
		{
			name:    "出错记为 failure",
			level:   logger.Info,
			elapsed: time.Millisecond,
			err:     errors.New("syntax error"),
			want:    kindFailure,
		},
		{
			name:    "出错优先于超时",
			level:   logger.Info,
			elapsed: time.Second,
			err:     errors.New("syntax error"),
			want:    kindFailure,
		},
		{
			name:    "记录未找到且不忽略时记为 failure",
			level:   logger.Info,
			elapsed: time.Millisecond,
			err:     logger.ErrRecordNotFound,
			want:    kindFailure,
		},
		{
			name:     "记录未找到且要求忽略时降级为普通查询",
			level:    logger.Info,
			elapsed:  time.Millisecond,
			err:      logger.ErrRecordNotFound,
			ignoreNF: true,
			want:     kindQuery,
		},
		{
			name:    "超过慢阈值记为 slow",
			level:   logger.Info,
			elapsed: time.Second,
			want:    kindSlow,
		},
		{
			name:            "慢阈值为 0 表示不判慢",
			level:           logger.Info,
			elapsed:         time.Second,
			noSlowThreshold: true,
			want:            kindQuery,
		},
		{
			name:    "普通查询记为 query",
			level:   logger.Info,
			elapsed: time.Millisecond,
			want:    kindQuery,
		},
		{
			name:    "Warn 级别不记普通查询",
			level:   logger.Warn,
			elapsed: time.Millisecond,
			want:    kindSkip,
		},
		{
			name:    "Error 级别不记普通查询",
			level:   logger.Error,
			elapsed: time.Millisecond,
			want:    kindSkip,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cfg := logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  tc.level,
				IgnoreRecordNotFoundError: tc.ignoreNF,
			}
			if tc.noSlowThreshold {
				cfg.SlowThreshold = 0
			}

			if got := classify(cfg, tc.elapsed, tc.err); got != tc.want {
				t.Errorf("classify() = %s, want %s", kindNames[got], kindNames[tc.want])
			}
		})
	}
}

// TestTraceSkipsSQLFormatting 固定「不记录就不做 SQL 格式化」这个取舍：
// fc 由 gorm 提供，展开 SQL 并不便宜，不该为不输出的查询付出。
func TestTraceSkipsSQLFormatting(t *testing.T) {
	called := false
	fc := func() (string, int64) {
		called = true
		return "SELECT 1", 1
	}

	tests := []struct {
		name   string
		level  logger.LogLevel
		wantFc bool
	}{
		{name: "Silent 级别不调用 fc", level: logger.Silent},
		{name: "Warn 级别下的普通查询不调用 fc", level: logger.Warn},
		{name: "Info 级别下的普通查询要调用 fc", level: logger.Info, wantFc: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			redirect(t)
			called = false

			New(testConfig(tc.level, false)).
				Trace(context.Background(), time.Now(), fc, nil)

			if called != tc.wantFc {
				t.Errorf("fc 是否被调用 = %v, want %v", called, tc.wantFc)
			}
		})
	}
}

// assertFields 逐项比对字段的键与值。
func assertFields(t *testing.T, got, want []logx.LogField) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("字段数 = %d, want %d：%+v", len(got), len(want), got)
	}
	for i, w := range want {
		if got[i].Key != w.Key || got[i].Value != w.Value {
			t.Errorf("字段[%d] = %s:%v, want %s:%v", i, got[i].Key, got[i].Value, w.Key, w.Value)
		}
	}
}

// TestQueryFields 固定字段的组装结果，含出错时多出的 gorm.error。
// duration 断言的是字面量，因为它是 go-zero timex.ReprOfDuration 的对外格式。
func TestQueryFields(t *testing.T) {
	const (
		querySQL = "SELECT 1"
		elapsed  = 3 * time.Millisecond
	)

	t.Run("无错误时三个字段", func(t *testing.T) {
		assertFields(t, queryFields(querySQL, 7, elapsed, nil), []logx.LogField{
			logx.Field("gorm.sql", querySQL),
			logx.Field("gorm.rows", int64(7)),
			logx.Field("gorm.duration", "3.0ms"),
		})
	})

	t.Run("有错误时多一个 gorm.error", func(t *testing.T) {
		assertFields(t, queryFields(querySQL, 7, elapsed, errors.New("syntax error")), []logx.LogField{
			logx.Field("gorm.sql", querySQL),
			logx.Field("gorm.rows", int64(7)),
			logx.Field("gorm.duration", "3.0ms"),
			logx.Field("gorm.error", "syntax error"),
		})
	})
}
