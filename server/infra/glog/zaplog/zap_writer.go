package zaplog

import (
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// 日志按大小切割和备份个数、文件有效期，一般选zap+lumberjack
// 日志按天或者小时切割，一般选zap+file-rotatelogs
// 日志输出

// 按时间分割日志文件 runtime/log/2022-11
// rotate库于2022.11.9开始不再更新
func NewFileWriter(cfg ZapConfig) *rotatelogs.RotateLogs {
	fileName := cfg.FileName
	var opts []rotatelogs.Option
	// 最新的日志文件保存地址
	//opts = append(opts, rotatelogs.WithLinkName(path.Join(cfg.Path, fileName)))
	switch cfg.Rotation {
	case "size":
		// 设置日志备份数量。不能和日志备份时间一起设置。当达到日志备份数量的限制时，将会自动删除最旧的备份日志文件，以保持指定数量的备份日志文件。
		opts = append(opts, rotatelogs.WithRotationCount(30))
		// 设置单个日志文件大小
		opts = append(opts, rotatelogs.WithRotationSize(cfg.MaxSize))
		break
	case "daily":
	default:
		// 轮转的时间将会根据本地系统的时区来计算
		opts = append(opts, rotatelogs.WithClock(rotatelogs.Local))
		// 日志留存时间
		opts = append(opts, rotatelogs.WithMaxAge(time.Duration(cfg.KeepDays)*24*time.Hour))
		// 按照天来切割日志 %Y-%m-%d-%H-%M-%S 只会计算 %Y-%m-%d-00-00-00
		opts = append(opts, rotatelogs.WithRotationTime(time.Hour*24))
	}

	fileWriter, err := rotatelogs.New(
		path.Join(cfg.Path, "%Y-%m-%d", fileName), // 历史日志保存的地址
		opts...,
	)
	if err != nil {
		return nil
	}
	//fileWriter.Rotate()
	return fileWriter
}
