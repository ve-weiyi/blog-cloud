package glog

import (
	"encoding/json"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog/zaplog"
)

var defaultLogger *Glogger

func init() {
	defaultLogger = NewGlogger(1, zaplog.NewConsoleConfig())
}

func Init(skit int, cfg zaplog.ZapConfig) {
	defaultLogger = NewGlogger(skit, cfg)
}

func Default() *Glogger {
	return defaultLogger
}

func Error(v ...interface{}) {

	defaultLogger.SugaredLogger.Error(v...)
}

func Warn(v ...interface{}) {

	defaultLogger.SugaredLogger.Warn(v...)
}

func Info(v ...interface{}) {

	defaultLogger.SugaredLogger.Info(v...)
}

func Debug(v ...interface{}) {

	defaultLogger.SugaredLogger.Debug(v...)
}

// print 风格
func Errorf(template string, args ...interface{}) {

	defaultLogger.SugaredLogger.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {

	defaultLogger.SugaredLogger.Warnf(template, args...)
}

func Infof(template string, args ...interface{}) {

	defaultLogger.SugaredLogger.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {

	defaultLogger.SugaredLogger.Debugf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {

	defaultLogger.SugaredLogger.Errorw(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {

	defaultLogger.SugaredLogger.Warnw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {

	defaultLogger.SugaredLogger.Infow(msg, keysAndValues...)
}

func Debugw(msg string, keysAndValues ...interface{}) {

	defaultLogger.SugaredLogger.Debugw(msg, keysAndValues...)
}

func Println(v ...interface{}) {

	defaultLogger.SugaredLogger.Info(v...)
}

func JsonIndent(v ...interface{}) {

	bytes, _ := json.MarshalIndent(v, "", " ")
	defaultLogger.SugaredLogger.Info("--->", string(bytes))
}
