package gormlogx

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"gorm.io/gorm/logger"
)

// Config 仅封装行号过滤相关配置，职责单一
type Config struct {
	Skip         int      // 额外跳过的调用栈深度
	SkipKeywords []string // 跳过的文件关键词列表
}

// 默认过滤配置
var defaultConfig = Config{
	Skip:         0,
	SkipKeywords: []string{"/gen/", "/gorm/"},
}

type GormWriter struct {
	logger logger.Writer
	config Config
}

func NewGormWriter(w logger.Writer, c Config) *GormWriter {
	loggerWriter := w
	if loggerWriter == nil {
		loggerWriter = log.New(os.Stdout, "\r\n", log.LstdFlags)
	}

	return &GormWriter{
		logger: loggerWriter,
		config: c,
	}
}

// Printf 实现logger.Writer接口
// message:%s \n[%.3fms] [rows:%v] %s
// data:[line,time,rows,sql]
func (w *GormWriter) Printf(message string, data ...interface{}) {
	if len(data) == 0 {
		w.logger.Printf(message, data...)
		return
	}
	// 拦截打印行号，使用业务行号
	data[0] = FileWithLineNum(w.config.Skip, w.config.SkipKeywords)
	w.logger.Printf(message, data...)
}

// FileWithLineNum 函数返回当前文件的文件名和行号
func FileWithLineNum(skip int, skipKeywords []string) string {
	// i=2通常来自 GORM 内部。所以把 i 的起始值设为 2。
	// 叠加额外skip值为了快速定位到业务文件，i<15限制最大深度。
	for i := 2 + skip; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 关键词匹配：判断是否跳过当前文件
		skipFile := false
		for _, kw := range skipKeywords {
			if strings.Contains(file, kw) {
				skipFile = true
				break
			}
		}

		if !skipFile {
			return fmt.Sprintf("%s:%d", file, line)
		}
	}
	return ""
}
