package initialize

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
)

type gormWriter struct {
	logger.Writer
	ZapLogger *zap.Logger
	UseZapLog bool
}

// NewWriter gormWriter 构造函数
// Author [SliverHorn](https://github.com/SliverHorn)
func NewWriter(w logger.Writer) *gormWriter {
	return &gormWriter{
		Writer:    w,
		ZapLogger: global.LOG.Logger(),
		UseZapLog: false,
	}
}

// Printf 格式化打印日志
// message:%s \n[%.3fms] [rows:%v] %s
// data:[line,time,rows,sql]
func (w *gormWriter) Printf(message string, data ...interface{}) {
	data[0] = FileWithLineNum()
	if w.UseZapLog {
		//w.Writer.Printf(message, data...)
		global.LOG.Info(fmt.Sprintf(message, data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}

var gormSourceDir string = "gorm"
var genSourceDir string = "gen"

// 拦截打印行号，过滤gen和gorm的行号
// FileWithLineNum return the file name and line number of the current file
func FileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 5; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.Contains(file, gormSourceDir) && !strings.Contains(file, genSourceDir) && !strings.Contains(file, "repository")) {
			return fmt.Sprintf("%v:%v", file, strconv.FormatInt(int64(line), 10))
		}
	}
	return ""
}
