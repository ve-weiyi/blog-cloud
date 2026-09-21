package mq

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// handleEmailMessage 解析邮件消息并交给业务处理。
func handleEmailMessage(handler EmailHandler, msg *mqx.Message) error {
	var event EmailMessageEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logx.Errorf("[mq] unmarshal EmailMessageEvent: %v", err)
		return nil
	}
	return handler(context.Background(), &event)
}

// handleSmsMessage 解析短信消息并交给业务处理。
func handleSmsMessage(handler SmsHandler, msg *mqx.Message) error {
	var event SmsMessageEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logx.Errorf("[mq] unmarshal SmsMessageEvent: %v", err)
		return nil
	}
	return handler(context.Background(), &event)
}

// handleInboxMessage 解析站内信投递事件并交给业务处理。
func handleInboxMessage(handler InboxHandler, msg *mqx.Message) error {
	var event InboxMessageEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logx.Errorf("[mq] unmarshal InboxMessageEvent: %v", err)
		return nil
	}
	return handler(context.Background(), &event)
}

// handleLoginEvent 解析登录日志事件并交给业务处理。
func handleLoginEvent(handler LoginHandler, msg *mqx.Message) error {
	var event LoginEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logx.Errorf("[mq] unmarshal LoginEvent: %v", err)
		return nil
	}
	return handler(context.Background(), &event)
}

// handleLogoutEvent 解析登出事件并交给业务处理。
func handleLogoutEvent(handler LogoutHandler, msg *mqx.Message) error {
	var event LogoutEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		logx.Errorf("[mq] unmarshal LogoutEvent: %v", err)
		return nil
	}
	return handler(context.Background(), &event)
}

// chain 为消息处理加上 panic 兜底与重试：先兜 panic 保证消费循环不退出，
// 再重试 3 次。重试耗尽后返回错误，由 mqx 按 AutoAck 语义丢弃该消息，
// 避免毒消息无限重投。
func chain(next mqx.Handler) mqx.Handler {
	return mqx.Chain(
		mqx.Recovery(&mqx.LogxLogger{Logger: logx.WithContext(context.Background())}),
		mqx.Retry(3),
	)(next)
}
