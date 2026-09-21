// Package notifyx 定义通知事件的跨进程契约：rpc 侧发布，各网关实例订阅后推给在线客户端。
//
// 事件只声明「哪条通知消息被发布/撤回」，**不携带收件人**——收件人由客户端按自身身份
// 拉取，避免网关为此展开全量收件人列表。
package notifyx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
)

// 事件类型
const (
	// EventNotice 通知消息已发布
	EventNotice = "notice"
	// EventNoticeRevoke 通知消息已撤回
	EventNoticeRevoke = "notice-revoke"
)

// Event 通知事件载荷，JSON 键名与 HTTP 响应契约保持一致（snake_case）。
//
// 字段变更需**同步三处**：本文件、`protocol/api/admin/notification.api` 的
// NotifyStreamEvent、`blog-admin/src/composables/sse/sseEvents.ts`。
type Event struct {
	Type        string `json:"event"`
	MessageId   int64  `json:"message_id,omitempty"`
	Title       string `json:"title,omitempty"`
	Category    string `json:"category,omitempty"`
	Level       string `json:"level,omitempty"`
	PublishedAt int64  `json:"published_at,omitempty"`
}

// Encode 序列化为可发布到 channel 的载荷
func (e Event) Encode() (string, error) {
	b, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Decode 解析 channel 载荷
func Decode(payload string) (Event, error) {
	var e Event
	err := json.Unmarshal([]byte(payload), &e)
	return e, err
}

// 重试参数：合计最坏约 150ms，不会明显拖慢调用方
const (
	publishAttempts = 3
	publishBackoff  = 50 * time.Millisecond
)

// Publish 把事件广播给所有网关实例。
//
// 失败会**有限重试**：一次 Redis 抖动就丢掉信号，会让客户端长时间停在旧未读数
// （前端虽有兜底对账，但那是分钟级的）。全部尝试仍失败时返回错误，
// 由调用方决定是否记日志——推送失败不影响发布/撤回本身。
func Publish(ctx context.Context, rdb redis.UniversalClient, e Event) error {
	payload, err := e.Encode()
	if err != nil {
		return err
	}

	var lastErr error
	backoff := publishBackoff
	for attempt := 1; attempt <= publishAttempts; attempt++ {
		if err := rdb.Publish(ctx, cachekey.NotifyEventChannel, payload).Err(); err != nil {
			lastErr = err
			if attempt < publishAttempts {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(backoff):
				}
				backoff *= 2
			}
			continue
		}
		return nil
	}
	return lastErr
}
