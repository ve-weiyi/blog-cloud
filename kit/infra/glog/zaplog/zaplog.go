package zaplog

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 默认的zap不需要增加skip
func NewZapLogger(skip int, cfg ZapConfig) *zap.Logger {
	// 使用了core的NewTee
	core := zapcore.NewTee(
		NewTreeCore(cfg)...,
	)
	// 创建一个将日志写入 WriteSyncer 的核心。
	// Glogger.Debug->skip1  glog.Debug->skip2
	logger := zap.New(
		core,
	)

	logger = logger.WithOptions(zap.AddCaller())
	logger = logger.WithOptions(zap.AddCallerSkip(skip))

	return logger
}

// 不同等级的日志使用不同的日志输出
func NewTreeCore(cfg ZapConfig) []zapcore.Core {
	//使用不同后缀收集日志
	var lvCores []zapcore.Core

	//输出到文件
	for level := cfg.TransportLevel(); level <= zapcore.FatalLevel; level++ {
		//if cfg.ShowColorPrefix() {
		//	cfg.Prefix = level.CapitalString() + " "
		//}
		encoder := NewEncoder(cfg)
		writer := zapcore.AddSync(NewWriter(cfg))
		lvCores = append(lvCores, zapcore.NewCore(encoder, writer, GetLevelPriority(level)))
	}

	return lvCores
}

func NewEncoder(cfg ZapConfig) zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	//encodeConfig := zap.NewProductionEncoderConfig()
	//encodeConfig := zap.NewDevelopmentEncoderConfig()
	encodeConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "log",
		TimeKey:        "time",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     cfg.PrefixTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	if cfg.ShowColor {
		encodeConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	if cfg.ShowLine {
		encodeConfig.EncodeCaller = zapcore.FullCallerEncoder
	} else {
		encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	switch cfg.Encoding {
	case FormatConsole: // 控制台格式
		return zapcore.NewConsoleEncoder(encodeConfig)
	case FormatJson: // json格式
		return zapcore.NewJSONEncoder(encodeConfig)
	}
	//json格式
	return zapcore.NewJSONEncoder(encodeConfig)
}

func NewWriter(cfg ZapConfig) zapcore.WriteSyncer {
	switch cfg.Mode {
	case "console":
		return zapcore.AddSync(os.Stderr)
	case "file":
		return zapcore.AddSync(NewFileWriter(cfg))
	default:
		return zapcore.AddSync(os.Stderr)
	}
}
