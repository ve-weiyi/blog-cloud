package mq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// ErrUnavailable 消息队列未就绪（RabbitMQ 不可用），消息未投递。
// 调用方据此区分「投递能力缺失」与「投递失败」，前者通常降级处理。
var ErrUnavailable = errors.New("mq: broker unavailable")

// publish 将事件序列化后投递到指定交换机。
// 发布器为 nil（RabbitMQ 不可用）时返回 ErrUnavailable。
func publish(ctx context.Context, pub mqx.Publisher, routingKey string, event any) error {
	if pub == nil {
		return ErrUnavailable
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("mq: marshal %s: %w", routingKey, err)
	}

	return pub.Publish(ctx, routingKey, &mqx.Message{Body: body})
}

// ── 邮件 ──

// PublishEmailMessageEvent 投递邮件消息事件。
func PublishEmailMessageEvent(ctx context.Context, event *EmailMessageEvent) error {
	return publish(ctx, emailPub, EmailRoutingKey, event)
}

// ── 短信 ──

// PublishSmsMessageEvent 投递短信消息事件。
func PublishSmsMessageEvent(ctx context.Context, event *SmsMessageEvent) error {
	return publish(ctx, smsPub, SmsRoutingKey, event)
}

// ── 站内信 ──

// PublishInboxMessageEvent 投递站内信投递事件。
func PublishInboxMessageEvent(ctx context.Context, event *InboxMessageEvent) error {
	return publish(ctx, inboxPub, InboxRoutingKey, event)
}

// ── 登录日志 ──

// PublishLoginEvent 投递登录日志事件。
func PublishLoginEvent(ctx context.Context, event *LoginEvent) error {
	return publish(ctx, loginPub, LoginRoutingKey, event)
}

// ── 登出 ──

// PublishLogoutEvent 投递登出事件。
func PublishLogoutEvent(ctx context.Context, event *LogoutEvent) error {
	return publish(ctx, logoutPub, LogoutRoutingKey, event)
}
