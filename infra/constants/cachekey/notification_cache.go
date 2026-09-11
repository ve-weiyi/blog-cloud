package cachekey

import "fmt"

// Notification模块缓存Key定义

const (
	// InboxReadPrefix 站内消息已读记录前缀（zset，member=messageId，score=readAt毫秒时间戳）
	// Key格式: notification:inbox:read:{userId}
	InboxReadPrefix = "notification:inbox:read:"
)

// GetInboxReadKey 获取用户站内消息已读zset的Key
func GetInboxReadKey(userId string) string {
	return fmt.Sprintf("%s%s", InboxReadPrefix, userId)
}
