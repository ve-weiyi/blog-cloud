package task

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/common/rediskey"
)

type TaskClearChatOnline struct {
	Redis *redis.Redis

	expires time.Duration
}

func NewTaskClearChatOnline(redis *redis.Redis) *TaskClearChatOnline {
	return &TaskClearChatOnline{
		Redis:   redis,
		expires: time.Minute * 5,
	}
}

func (t *TaskClearChatOnline) Run() {
	redisClient := t.Redis

	expiredTime := time.Now().Add(-t.expires).Unix()
	onlineKey := rediskey.GetChatOnlineKey()

	logger := NewLogger()

	count, err := redisClient.ZcardCtx(context.Background(), onlineKey)
	if err != nil {
		logger.AddLog("Error getting online count: %v", err)
	}

	if count > 0 {
		// 移除所有过期元素
		removedCount, err := redisClient.ZremrangebyscoreCtx(context.Background(), onlineKey, 0, expiredTime)
		if err != nil {
			logger.AddLog("Error removing expired elements: %v", err)
		}
		if removedCount > 0 {
			logger.AddLog("Removed %d expired elements", removedCount)
		}
	}
}

type Logger struct {
	Logs map[int64]string
}

func NewLogger() *Logger {
	return &Logger{
		Logs: make(map[int64]string),
	}
}

func (l *Logger) AddLog(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Logs[time.Now().Unix()] = msg
	fmt.Println(msg)
}
