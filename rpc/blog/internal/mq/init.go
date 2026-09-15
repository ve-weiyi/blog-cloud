package mq

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/config"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/adapter/mqx/rabbitmqx"
)

var (
	emailMQ  mqx.MessageQueue
	smsMQ    mqx.MessageQueue
	inboxMQ  mqx.MessageQueue
	loginMQ  mqx.MessageQueue
	logoutMQ mqx.MessageQueue

	emailPub  mqx.Publisher
	smsPub    mqx.Publisher
	inboxPub  mqx.Publisher
	loginPub  mqx.Publisher
	logoutPub mqx.Publisher

	emailSub  mqx.Subscriber
	smsSub    mqx.Subscriber
	inboxSub  mqx.Subscriber
	loginSub  mqx.Subscriber
	logoutSub mqx.Subscriber
)

// Init 建立各交换机的连接，并派生发布器与订阅器，供全局复用。
// 订阅器不在此启动，由 StartXxxSubscriber 显式启动。
// RabbitMQ 不可用时只记录错误并降级：发布器保持 nil，投递返回错误而不阻断主流程。
func Init(rc config.RabbitMQConf) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", rc.Username, rc.Password, rc.Host, rc.Port)

	emailMQ = newMQ(url, EmailExchange)
	smsMQ = newMQ(url, SmsExchange)
	inboxMQ = newMQ(url, InboxExchange)
	loginMQ = newMQ(url, LoginExchange)
	logoutMQ = newMQ(url, LogoutExchange)

	emailPub = newPublisher(emailMQ, EmailExchange)
	smsPub = newPublisher(smsMQ, SmsExchange)
	inboxPub = newPublisher(inboxMQ, InboxExchange)
	loginPub = newPublisher(loginMQ, LoginExchange)
	logoutPub = newPublisher(logoutMQ, LogoutExchange)

	emailSub = newSubscriber(emailMQ, EmailExchange)
	smsSub = newSubscriber(smsMQ, SmsExchange)
	inboxSub = newSubscriber(inboxMQ, InboxExchange)
	loginSub = newSubscriber(loginMQ, LoginExchange)
	logoutSub = newSubscriber(logoutMQ, LogoutExchange)
}

// Close 关闭全部订阅器、发布器与连接。
func Close() {
	for _, sub := range []mqx.Subscriber{emailSub, smsSub, inboxSub, loginSub, logoutSub} {
		if sub != nil {
			_ = sub.Close()
		}
	}
	for _, pub := range []mqx.Publisher{emailPub, smsPub, inboxPub, loginPub, logoutPub} {
		if pub != nil {
			_ = pub.Close()
		}
	}
	for _, queue := range []mqx.MessageQueue{emailMQ, smsMQ, inboxMQ, loginMQ, logoutMQ} {
		if queue != nil {
			_ = queue.Close()
		}
	}
}

// newMQ 创建消息队列连接。失败时返回 nil，表示该交换机的收发能力不可用。
func newMQ(url, exchange string) mqx.MessageQueue {
	queue, err := rabbitmqx.New(&rabbitmqx.Config{
		URL:          url,
		ExchangeName: exchange,
		ExchangeType: rabbitmqx.ExchangeFanout,
		Durable:      true,
		AutoDelete:   false,
		Logger:       &mqx.LogxLogger{Logger: logx.WithContext(context.Background())},
	})
	if err != nil {
		logx.Errorf("[mq] init %s: %v", exchange, err)
		return nil
	}
	return queue
}

// newPublisher 从连接派生发布器。
func newPublisher(queue mqx.MessageQueue, exchange string) mqx.Publisher {
	if queue == nil {
		return nil
	}

	pub, err := queue.Publisher()
	if err != nil {
		logx.Errorf("[mq] create publisher %s: %v", exchange, err)
		return nil
	}
	return pub
}

// newSubscriber 从连接派生订阅器。
func newSubscriber(queue mqx.MessageQueue, exchange string) mqx.Subscriber {
	if queue == nil {
		return nil
	}

	sub, err := queue.Subscriber()
	if err != nil {
		logx.Errorf("[mq] create subscriber %s: %v", exchange, err)
		return nil
	}
	return sub
}

// ── 消费者启动 ──

// EmailHandler 处理邮件消息事件。
type EmailHandler func(ctx context.Context, event *EmailMessageEvent) error

// SmsHandler 处理短信消息事件。
type SmsHandler func(ctx context.Context, event *SmsMessageEvent) error

// InboxHandler 处理站内信投递事件。
type InboxHandler func(ctx context.Context, event *InboxMessageEvent) error

// LoginHandler 处理登录日志事件。
type LoginHandler func(ctx context.Context, event *LoginEvent) error

// LogoutHandler 处理登出事件。
type LogoutHandler func(ctx context.Context, event *LogoutEvent) error

// StartEmailSubscriber 启动邮件消息消费者（goroutine）。
func StartEmailSubscriber(handler EmailHandler) {
	go start(EmailExchange, emailSub, EmailRoutingKey, EmailQueue, func(ctx context.Context, msg *mqx.Message) error {
		return handleEmailMessage(handler, msg)
	})
}

// StartSmsSubscriber 启动短信消息消费者（goroutine）。
func StartSmsSubscriber(handler SmsHandler) {
	go start(SmsExchange, smsSub, SmsRoutingKey, SmsQueue, func(ctx context.Context, msg *mqx.Message) error {
		return handleSmsMessage(handler, msg)
	})
}

// StartInboxSubscriber 启动站内信投递消费者（goroutine）。
func StartInboxSubscriber(handler InboxHandler) {
	go start(InboxExchange, inboxSub, InboxRoutingKey, InboxQueue, func(ctx context.Context, msg *mqx.Message) error {
		return handleInboxMessage(handler, msg)
	})
}

// StartLoginSubscriber 启动登录日志消费者（goroutine）。
func StartLoginSubscriber(handler LoginHandler) {
	go start(LoginExchange, loginSub, LoginRoutingKey, LoginQueue, func(ctx context.Context, msg *mqx.Message) error {
		return handleLoginEvent(handler, msg)
	})
}

// StartLogoutSubscriber 启动登出事件消费者（goroutine）。
func StartLogoutSubscriber(handler LogoutHandler) {
	go start(LogoutExchange, logoutSub, LogoutRoutingKey, LogoutQueue, func(ctx context.Context, msg *mqx.Message) error {
		return handleLogoutEvent(handler, msg)
	})
}

// start 在独立 goroutine 中运行订阅循环。
// 订阅建立失败只记录日志，不中断服务：其余交换机与 RPC 服务不受影响。
func start(exchange string, sub mqx.Subscriber, routingKey, queue string, handler mqx.Handler) {
	if sub == nil {
		logx.Errorf("[mq] %s subscriber unavailable, consume disabled", exchange)
		return
	}

	logx.Infof("[mq] %s subscriber starting", exchange)
	if err := sub.Subscribe(context.Background(), routingKey, chain(handler),
		mqx.WithGroup(queue),
		mqx.WithPrefetch(10),
	); err != nil {
		logx.Errorf("[mq] %s subscriber stopped: %v", exchange, err)
	}
}
