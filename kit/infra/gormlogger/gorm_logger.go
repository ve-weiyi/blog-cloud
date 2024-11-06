package gormlogger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gorm.io/gorm/logger"
)

type Option func(w *GormWriter)

func AddSkip(skip int) Option {
	return func(w *GormWriter) {
		w.addSkip = skip
	}
}

type GormWriter struct {
	addSkip int
	l       logger.Writer
}

func NewGormWriter(opts ...Option) *GormWriter {
	w := &GormWriter{
		l: log.New(os.Stdout, "\r\n", log.LstdFlags),
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// Printf 格式化打印日志
// message:%s \n[%.3fms] [rows:%v] %s
// data:[line,time,rows,sql]
func (w *GormWriter) Printf(message string, data ...interface{}) {
	data[0] = w.FileWithLineNum()
	w.l.Printf(message, data...)
}

var gormSourceDir string = "gorm"
var genSourceDir string = "gen"

// 拦截打印行号，过滤gen和gorm的行号
// FileWithLineNum return the file name and line number of the current file
func (w *GormWriter) FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 5; i < 15; i++ {
		_, ff, _, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(ff, gormSourceDir)) {
			_, file, line, ok := runtime.Caller(i + w.addSkip)
			if ok {
				return fmt.Sprintf("%v:%v", file, strconv.FormatInt(int64(line), 10))
			}
		}
	}
	return ""
}
