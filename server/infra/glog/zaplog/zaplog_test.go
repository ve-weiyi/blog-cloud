package zaplog

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZap(t *testing.T) {

	// 配置日志编码器
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cfg.Build()

	// 输出彩色日志
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

}

// 输出到控制台测试
func TestConsole(t *testing.T) {
	cfg := NewConsoleConfig()

	logc := NewZapLogger(0, cfg)
	logc.Sugar().Debugw("debug", zap.String("key", "value"))
	logc.Sugar().Info("info")
	logc.Sugar().Warn("warn")
	logc.Sugar().Error("error")
	//logc.Sugar().DPanic("dpanic")
	//logc.Sugar().Panic("panic")
	//logc.Sugar().Fatal("fatal")
}

// 输出到文件测试
func TestFile(t *testing.T) {
	cfg := NewFileConfig()
	//cfg.Encoding = FormatConsole
	cfg.ShowColor = true

	logc := NewZapLogger(0, cfg)
	logc.Sugar().Debugw("debug", zap.String("key", "value"))
	logc.Sugar().Info("info")
	logc.Sugar().Warn("warn")
	logc.Sugar().Error("error")
	//logc.Sugar().DPanic("dpanic")
	//logc.Sugar().Panic("panic")
	//logc.Sugar().Fatal("fatal")
}

// 同时输出到控制台和文件测试
func TestConsoleAndFile(t *testing.T) {
	cfg := NewFileConfig()
	cfg.Encoding = FormatConsole
	logs := NewZapLogger(0, cfg)

	// 压测

	logs.Sugar().Info("hello")
}

// 压力测试
func BenchmarkExample(b *testing.B) {
	logs := NewZapLogger(0, NewConsoleConfig())
	logj := NewZapLogger(0, NewFileConfig())
	for i := 0; i < b.N; i++ {
		// 在这里执行需要进行压力测试的代码
		logs.Sugar().Info("hello", i)
		logj.Sugar().Info("hello", i)
	}
}

func TestWriter(t *testing.T) {
	writer := os.Stderr

	write, err := writer.Write([]byte("hello\n"))
	fmt.Println("write", write, err)
	if err != nil {
		return
	}
}
