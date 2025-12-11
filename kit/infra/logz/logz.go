package logz

import (
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志级别常量
const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type LogConfig struct {
	Level      string `json:"level"`
	Mode       string `json:"mode"` // console|file
	Filename   string `json:"file_name"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

// SetLog 创建新的日志实例
func SetLog(conf *LogConfig, opts ...zap.Option) *zap.Logger {
	// 自动创建日志目录
	logDir := filepath.Dir(conf.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		// 如果创建目录失败，使用当前目录
		conf.Filename = filepath.Base(conf.Filename)
	}

	// 配置 lumberjack 日志轮转
	hook := lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		Compress:   conf.Compress,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch conf.Level {
	case LogLevelDebug:
		atomicLevel.SetLevel(zap.DebugLevel)
	case LogLevelInfo:
		atomicLevel.SetLevel(zap.InfoLevel)
	case LogLevelWarn:
		atomicLevel.SetLevel(zap.WarnLevel)
	case LogLevelError:
		atomicLevel.SetLevel(zap.ErrorLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	var cores []zapcore.Core

	if strings.Contains(conf.Mode, "console") {
		// 文件输出使用JSON格式
		consoleEncoder := zapcore.NewConsoleEncoder(getConsoleEncoderConfig())
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), atomicLevel))
	}

	if strings.Contains(conf.Mode, "file") {
		// 控制台输出使用带颜色的文本格式
		fileEncoder := zapcore.NewJSONEncoder(getJsonEncoderConfig())
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(&hook), atomicLevel))
	}

	core := zapcore.NewTee(
		cores...,
	)

	// 创建日志器
	defaultLogger = zap.New(core, append(opts, zap.AddCaller(), zap.AddCaller())...)
	return defaultLogger
}

// getJsonEncoderConfig 获取JSON编码器配置
func getJsonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// getConsoleEncoderConfig 获取控制台编码器配置（带颜色）
func getConsoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,        // 带颜色的级别
		EncodeTime:     zapcore.TimeEncoderOfLayout("15:04:05"), // 简化时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 简化文件路径
		EncodeName:     zapcore.FullNameEncoder,
	}
}

// go使用zap + lumberjack重构项目的日志系统 https://blog.csdn.net/weixin_52000204/article/details/126651319
var defaultLogger *zap.Logger

func init() {
	defaultLogger = SetLog(&LogConfig{
		Level:      "info",
		Mode:       "console|file",
		Filename:   "./app.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   false,
	})
}

func L() *zap.Logger {
	return defaultLogger
}

func S() *zap.SugaredLogger {
	return defaultLogger.Sugar()
}
