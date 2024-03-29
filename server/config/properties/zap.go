package properties

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

type Zap struct {
	ServerName     string `mapstructure:"server-name" json:"server-name" yaml:"server-name"`             // 服务名称
	Mode           string `mapstructure:"mode" json:"mode" yaml:"mode"`                                  // 模式
	Format         string `mapstructure:"format" json:"format" yaml:"format"`                            // 输出
	Level          string `mapstructure:"level" json:"level" yaml:"level"`                               // 级别
	Prefix         string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                            // 日志前缀
	EncodeLevel    string `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`          // 编码级
	EncodeCaller   string `mapstructure:"encode-caller" json:"encode-caller" yaml:"encode-caller"`       // 编码调用者
	EncodeColorful bool   `mapstructure:"encode-colorful" json:"encode-colorful" yaml:"encode-colorful"` // 编码调用者

	CacheDir string `mapstructure:"cache-dir" json:"cache-dir" yaml:"cache-dir"` // 日志文件夹
	MaxAge   int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`       // 日志留存时间
}

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *Zap) ZapEncodeLevel() zapcore.LevelEncoder {
	switch {
	case z.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case z.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case z.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case z.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 根据字符串转化为 zapcore.Level
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *Zap) TransportLevel() zapcore.Level {
	z.Level = strings.ToLower(z.Level)
	switch z.Level {
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
