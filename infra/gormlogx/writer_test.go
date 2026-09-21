package gormlogx

import (
	"slices"
	"strings"
	"testing"
)

// recordWriter 实现 gorm 的 logger.Writer，记下最后一次调用。
type recordWriter struct {
	calls int
	msg   string
	args  []any
}

func (w *recordWriter) Printf(msg string, args ...any) {
	w.calls++
	w.msg = msg
	w.args = args
}

func TestNewGormWriter(t *testing.T) {
	t.Run("未注入 Writer 时回落标准输出", func(t *testing.T) {
		if w := NewGormWriter(nil, Config{}); w.logger == nil {
			t.Fatal("应回落为一个非空的 logger.Writer")
		}
	})

	t.Run("未配置跳过关键词时由默认值补齐", func(t *testing.T) {
		w := NewGormWriter(&recordWriter{}, Config{})

		want := []string{"/gen/", "/gorm/"}
		if !slices.Equal(w.config.SkipKeywords, want) {
			t.Errorf("SkipKeywords = %v, want %v", w.config.SkipKeywords, want)
		}
	})

	t.Run("已配置跳过关键词时不被默认值覆盖", func(t *testing.T) {
		w := NewGormWriter(&recordWriter{}, Config{SkipKeywords: []string{"/x/"}})

		if want := []string{"/x/"}; !slices.Equal(w.config.SkipKeywords, want) {
			t.Errorf("SkipKeywords = %v, want %v", w.config.SkipKeywords, want)
		}
	})
}

func TestGormWriterPrintf(t *testing.T) {
	t.Run("无参数时原样透传", func(t *testing.T) {
		inner := &recordWriter{}
		NewGormWriter(inner, Config{}).Printf("plain message")

		if inner.calls != 1 {
			t.Errorf("调用次数 = %d, want 1", inner.calls)
		}
		if inner.msg != "plain message" {
			t.Errorf("文案 = %q, want %q", inner.msg, "plain message")
		}
		if len(inner.args) != 0 {
			t.Errorf("参数 = %v, want 空", inner.args)
		}
	})

	t.Run("有参数时把 data[0] 改写成业务行号", func(t *testing.T) {
		inner := &recordWriter{}
		// 关键词命中栈上所有帧，行号必为空串 —— 便于确定性地断言「data[0] 被改写」这件事本身。
		NewGormWriter(inner, Config{SkipKeywords: []string{"/"}}).
			Printf("%s [%.3fms] [rows:%v] %s", "gorm 内部行号", 1.5, 1, "SELECT 1")

		if inner.calls != 1 {
			t.Fatalf("调用次数 = %d, want 1", inner.calls)
		}
		if got := inner.args[0]; got != "" {
			t.Errorf("data[0] = %v, want 空串（被业务行号覆盖）", got)
		}
	})
}

// callFileLine 只是多包一层，让 fileLine 的调用栈里多一帧本文件，
// 这样 skip=0 时它能定位到本测试文件而不是 testing 包的内部帧。
func callFileLine(skip int, keywords []string) string {
	return fileLine(skip, keywords)
}

func TestFileLine(t *testing.T) {
	const selfFile = "writer_test.go"

	tests := []struct {
		name         string
		keywords     []string
		wantEmpty    bool
		wantContains string
		wantOmit     string
	}{
		{
			name:         "无跳过关键词时定位到调用方文件",
			wantContains: selfFile,
		},
		{
			name:     "关键词命中调用方文件时跳过它",
			keywords: []string{selfFile},
			wantOmit: selfFile,
		},
		{
			name:      "所有帧都被关键词命中时返回空串",
			keywords:  []string{"/"},
			wantEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := callFileLine(0, tc.keywords)

			if tc.wantEmpty {
				if got != "" {
					t.Fatalf("got %q, want 空串", got)
				}
				return
			}
			if got == "" {
				t.Fatal("返回空串，没有定位到任何帧")
			}
			if tc.wantContains != "" && !strings.Contains(got, tc.wantContains) {
				t.Errorf("got %q, want 含 %q", got, tc.wantContains)
			}
			if tc.wantOmit != "" && strings.Contains(got, tc.wantOmit) {
				t.Errorf("got %q, want 不含 %q", got, tc.wantOmit)
			}
		})
	}
}
