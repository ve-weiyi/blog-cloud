package mq

import (
	"context"
	"log/slog"

	"github.com/zeromicro/go-zero/core/logx"
)

// newSlogLogger 把 slog 输出桥接到 go-zero 的 logx，保持日志格式与配置一致。
func newSlogLogger() *slog.Logger {
	return slog.New(&logxHandler{})
}

type logxHandler struct {
	attrs []slog.Attr
}

func (h *logxHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *logxHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := make([]logx.LogField, 0, len(h.attrs)+record.NumAttrs())
	for _, attr := range h.attrs {
		fields = append(fields, logx.Field(attr.Key, attr.Value.Any()))
	}
	record.Attrs(func(attr slog.Attr) bool {
		fields = append(fields, logx.Field(attr.Key, attr.Value.Any()))
		return true
	})

	logger := logx.WithContext(ctx)
	if record.Level >= slog.LevelError {
		logger.Errorw(record.Message, fields...)
	} else {
		logger.Infow(record.Message, fields...)
	}
	return nil
}

func (h *logxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	combined := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	combined = append(combined, h.attrs...)
	combined = append(combined, attrs...)
	return &logxHandler{attrs: combined}
}

func (h *logxHandler) WithGroup(string) slog.Handler { return h }
