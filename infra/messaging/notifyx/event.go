// Package notifyx 定义通知事件的跨进程契约：rpc 侧发布，各网关实例订阅后推给在线客户端。
//
// 事件只声明「哪条通知消息变了」：管理端的公告事件广播到共享频道，读者端的收件箱信号
// 按收件人投递到个人频道。事件始终不带正文与收件人列表 —— 客户端收到后按自身身份回查，
// 避免网关为此展开全量收件人。
package notifyx

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
)

// 事件类型。管理端只关心公告状态变化；读者端的「收件箱变了」是另一条流。
const (
	// EventAnnouncementPublished 公告已发布
	EventAnnouncementPublished = "announcement-published"
	// EventAnnouncementRevoked 公告已撤回
	EventAnnouncementRevoked = "announcement-revoked"
	// EventInboxChanged 收件人的收件箱发生了变化
	EventInboxChanged = "notify"
)

// Event 通知事件载荷，JSON 键名与 HTTP 响应契约保持一致（snake_case）。
//
// 字段变更需同步全部订阅端：本文件、`protocol/api/admin/notification.api` 的
// AnnouncementStreamEvent、`blog-admin/src/composables/sse/sseEvents.ts`。
type Event struct {
	Type           string `json:"event"`
	AnnouncementId int64  `json:"announcement_id,omitempty"`
	Title          string `json:"title,omitempty"`
	Category       string `json:"category,omitempty"`
	Level          string `json:"level,omitempty"`
	PublishedAt    int64  `json:"published_at,omitempty"`
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
	return publishWithRetry(ctx, rdb, cachekey.NotifyEventChannel, payload)
}

// PublishToUser 向某个收件人的频道发一条「收件箱变了」的信号。
//
// 信号不带正文，也不带收件人以外的信息：客户端收到后回查拉取接口，
// 离线期间与丢信号都由定时对账兜底。
func PublishToUser(ctx context.Context, rdb redis.UniversalClient, userID string) error {
	if userID == "" {
		return nil
	}
	payload, err := Event{Type: EventInboxChanged}.Encode()
	if err != nil {
		return err
	}
	return publishWithRetry(ctx, rdb, cachekey.NotifyUserChannel(userID), payload)
}

func publishWithRetry(ctx context.Context, rdb redis.UniversalClient, channel, payload string) error {
	var lastErr error
	backoff := publishBackoff
	for attempt := 1; attempt <= publishAttempts; attempt++ {
		if err := rdb.Publish(ctx, channel, payload).Err(); err != nil {
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
