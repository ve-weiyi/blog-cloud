package zaplog

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	FormatConsole = "console"
	FormatJson    = "json"
)

// zapcore.Core 需要三个配置
// ——Encoder
// ——WriteSyncer
// ——LogLevel

type ZapConfig struct {
	// 日志输出位置
	Mode     string `json:",default=console,options=[console,file,volume]"`
	Encoding string `json:"encoding"` // 输出格式 使用json还是console

	// 日志输出格式设置
	Level     string `json:"level" `     // 级别,低于该级别时不打印日志
	Prefix    string `json:"prefix" `    // 日志前缀
	ShowLine  bool   `json:"show-line"`  // 显示行号
	ShowColor bool   `json:"show-color"` // 显示颜色

	// 日志保存文件设置
	Path       string `json:"file-dir"`    //日志文件的路径
	FileName   string `json:"file-name"`   // 日志文件名
	Compress   bool   `json:"compress"`    //是否压缩/归档旧文件
	KeepDays   int    `json:"keep_days"`   //保留旧文件的最大天数。默认30天
	MaxBackups int    `json:"max_backups"` //保留旧文件的最大个数。MaxAge和MaxBackups只能同时使用一个
	MaxSize    int64  `json:"max_size"`    //在进行切割之前，日志文件的最大大小（以MB为单位）
	Rotation   string `json:"rotation,default=daily,options=[daily,size]"`
}

func NewConsoleConfig() ZapConfig {
	// 返回默认的配置
	return ZapConfig{
		Mode:       "console",
		Level:      "debug",
		Prefix:     "",
		Encoding:   FormatConsole,
		ShowLine:   true,
		ShowColor:  true,
		Path:       fmt.Sprintf("runtime/logs"),
		MaxSize:    10 * 1024 * 1024, // 10MB
		KeepDays:   30,
		MaxBackups: 30,
		FileName:   "zap.log",
	}
}

func NewFileConfig() ZapConfig {
	// 返回默认的配置
	return ZapConfig{
		Mode:       "file",
		Level:      "debug",
		Prefix:     "",
		Encoding:   FormatJson,
		ShowLine:   false,
		ShowColor:  false,
		FileName:   "zap.log",
		Path:       fmt.Sprintf("runtime/logs"),
		MaxSize:    10 * 1024 * 1024, // 10MB
		KeepDays:   30,
		MaxBackups: 30,
	}
}

// 编码器,可选项 LowercaseLevelEncoder LowercaseColorLevelEncoder CapitalLevelEncoder CapitalColorLevelEncoder
// ZapEncodeLevel 根据 LevelEncoder 返回 zapcore.LevelEncoder
func (m *ZapConfig) ZapEncodeLevel() zapcore.LevelEncoder {
	if m.ShowColor {
		return zapcore.CapitalColorLevelEncoder // 大写编码器带颜色
	}

	return zapcore.CapitalLevelEncoder // 大写编码器
}

// TransportLevel 根据字符串转化为 zapcore.Level
func (m *ZapConfig) TransportLevel() zapcore.Level {
	m.Level = strings.ToLower(m.Level)
	switch m.Level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// PrefixTimeEncoder 自定义带前缀日志输出时间格式
func (m *ZapConfig) PrefixTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(fmt.Sprintf("%s", m.Prefix) + t.Format("2006/01/02-15:04:05.000"))
}
