// Copyright © 2015-2018 Anker Innovations Technology Limited All Rights Reserved.
package glog

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog/zaplog"
)

// go使用zap + lumberjack重构项目的日志系统 https://blog.csdn.net/weixin_52000204/article/details/126651319
type Glogger struct {
	cfg                zaplog.ZapConfig // 日志配置
	log                *zap.Logger      // 并重性能与易用性，支持结构化和 printf 风格的日志记录。
	*zap.SugaredLogger                  // 非常强调性能，不提供 printf 风格的 api（减少了 interface{} 与 反射的性能损耗）
}

func NewGlogger(skip int, cfg zaplog.ZapConfig) *Glogger {
	logger := zaplog.NewZapLogger(skip, cfg)

	//初始化内部类
	mlog := new(Glogger)

	mlog.log = logger
	mlog.SugaredLogger = logger.Sugar()
	return mlog
}

func (mlog *Glogger) Logger() *zap.Logger {
	return mlog.log
}

func (mlog *Glogger) Error(v ...interface{}) {

	mlog.SugaredLogger.Error(v...)
}

func (mlog *Glogger) Warn(v ...interface{}) {

	mlog.SugaredLogger.Warn(v...)
}

func (mlog *Glogger) Info(v ...interface{}) {

	mlog.SugaredLogger.Info(v...)
}

func (mlog *Glogger) Debug(v ...interface{}) {

	mlog.SugaredLogger.Debug(v...)
}

// print 风格
func (mlog *Glogger) Errorf(template string, args ...interface{}) {

	mlog.SugaredLogger.Errorf(template, args...)
}

func (mlog *Glogger) Warnf(template string, args ...interface{}) {

	mlog.SugaredLogger.Warnf(template, args...)
}

func (mlog *Glogger) Infof(template string, args ...interface{}) {

	mlog.SugaredLogger.Infof(template, args...)
}

func (mlog *Glogger) Debugf(template string, args ...interface{}) {

	mlog.SugaredLogger.Debugf(template, args...)
}

func (mlog *Glogger) Errorw(msg string, keysAndValues ...interface{}) {

	mlog.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (mlog *Glogger) Warnw(msg string, keysAndValues ...interface{}) {

	mlog.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (mlog *Glogger) Infow(msg string, keysAndValues ...interface{}) {

	mlog.SugaredLogger.Infow(msg, keysAndValues...)
}

func (mlog *Glogger) Debugw(msg string, keysAndValues ...interface{}) {

	mlog.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (mlog *Glogger) Println(v ...interface{}) {

	mlog.SugaredLogger.Info(v...)
}

func (mlog *Glogger) JsonIndent(v ...interface{}) {

	bytes, _ := json.MarshalIndent(v, "", " ")
	mlog.SugaredLogger.Info("--->", string(bytes))
}
