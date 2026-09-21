package gormlogx

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"gorm.io/gorm/logger"
)

// Config 封装行号改写相关的配置。
type Config struct {
	Skip         int      // 额外跳过的调用栈深度
	SkipKeywords []string // 跳过的文件关键词列表
}

// defaultConfig 是未显式配置时使用的跳过规则，
// 目标是让行号落在业务代码上，而不是 gorm 自身的内部帧。
var defaultConfig = Config{
	Skip:         0,
	SkipKeywords: []string{"/gen/", "/gorm/"},
}

// GormWriter 实现 gorm 的 logger.Writer，把 gorm 打印的行号改写成业务行号。
type GormWriter struct {
	logger logger.Writer
	config Config
}

var _ logger.Writer = (*GormWriter)(nil)

// NewGormWriter 包装一个 gorm 的 Writer，w 为 nil 时回落到标准输出。
//
// 调用方没给跳过关键词时由 defaultConfig 补齐，因此传 Config{} 也能拿到默认行为，
// 不会退化成「不跳过任何帧」。
func NewGormWriter(w logger.Writer, c Config) *GormWriter {
	if w == nil {
		w = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}
	if len(c.SkipKeywords) == 0 {
		c.SkipKeywords = defaultConfig.SkipKeywords
	}

	return &GormWriter{
		logger: w,
		config: c,
	}
}

// Printf 实现 logger.Writer 接口。
// message:%s \n[%.3fms] [rows:%v] %s
// data:[line,time,rows,sql]
func (w *GormWriter) Printf(message string, data ...interface{}) {
	if len(data) == 0 {
		w.logger.Printf(message, data...)
		return
	}

	// 拦截打印行号，使用业务行号
	data[0] = fileLine(w.config.Skip, w.config.SkipKeywords)
	w.logger.Printf(message, data...)
}

// fileLine 返回调用栈上第一个未被跳过关键词命中的文件:行号，找不到时返回空串。
//
// i 从 2 起：帧 0、1 是 fileLine 自身与 gorm 的 Logger 包装，通常来自 GORM 内部，
// 所以把起始值设为 2。叠加额外 skip 值是为了快速定位到业务文件，i < 15 限制最大深度。
func fileLine(skip int, skipKeywords []string) string {
	for i := 2 + skip; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if !containsAny(file, skipKeywords) {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}

	return ""
}

// containsAny 判断 s 是否包含关键词列表中的任意一项。
func containsAny(s string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(s, kw) {
			return true
		}
	}

	return false
}
