package mq

import (
	"time"

	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/adapter/mqx/event"
)

// 交换机按业务域划分，routing key 即事件类型。
// 交换机与队列名是部署契约：改名必须按迁移处理，不能直接替换。
const (
	NotificationExchange = "blog.notification"
	AuditExchange        = "blog.audit"
	InteractionExchange  = "blog.interaction"
	NotificationDLQ      = "blog.notification.dlq"
	InteractionDLQ       = "blog.interaction.dlq"

	EmailRoutingKey  = "notification.email"
	SmsRoutingKey    = "notification.sms"
	InboxRoutingKey  = "notification.inbox"
	LoginRoutingKey  = "audit.login"
	LogoutRoutingKey = "audit.logout"
	LikeRoutingKey   = "interaction.like"

	EmailQueue       = "blog.rpc.notification.email"
	SmsQueue         = "blog.rpc.notification.sms"
	InboxQueue       = "blog.rpc.notification.inbox"
	AuditQueue       = "blog.rpc.audit"
	InteractionQueue = "blog.rpc.interaction"

	producer = "blog-rpc"
)

// LoginEvent 登录日志事件。
type LoginEvent struct {
	UserId     string `json:"user_id"`
	DeviceId   string `json:"device_id"`
	LoginType  string `json:"login_type"`
	Status     int64  `json:"status"`
	FailReason string `json:"fail_reason,omitempty"`
}

// LogoutEvent 登出事件。
type LogoutEvent struct {
	UserId     string `json:"user_id"`
	DeviceId   string `json:"device_id"`
	LogoutType string `json:"logout_type"` // manual-手动登出 timeout-超时登出 force-强制登出
}

// EmailMessageEvent 邮件消息事件。
type EmailMessageEvent struct {
	Email   string            `json:"email"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Scene   string            `json:"scene"`
	BizId   string            `json:"biz_id"`
	Params  map[string]string `json:"params,omitempty"`
}

// SmsMessageEvent 短信消息事件。
type SmsMessageEvent struct {
	Mobile string            `json:"mobile"`
	Scene  string            `json:"scene"`
	BizId  string            `json:"biz_id"`
	Params map[string]string `json:"params,omitempty"`
}

// 点赞事件携带的动作值。
const (
	LikeActionLike   = "like"
	LikeActionUnlike = "unlike"
)

// LikeEvent 点赞事件：动作在发布端按实时状态判定，消费端只做幂等收敛。
type LikeEvent struct {
	UserId     string `json:"user_id"`
	TargetType int64  `json:"target_type"`
	TargetId   int64  `json:"target_id"`
	Action     string `json:"action"` // like | unlike
}

// InboxMessageEvent 站内信投递事件：消费后为每个目标用户创建 delivery 记录。
type InboxMessageEvent struct {
	MessageId int64 `json:"message_id"`
}

// 事件类型与目的地、routing key 的绑定。发布与订阅统一使用这些值，
// 类型与绑定不匹配会在编译期报错。
var (
	EmailMessage = event.NewEventSpec[EmailMessageEvent](
		NotificationExchange, EmailRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
	SmsMessage = event.NewEventSpec[SmsMessageEvent](
		NotificationExchange, SmsRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
	InboxMessage = event.NewEventSpec[InboxMessageEvent](
		NotificationExchange, InboxRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
	LoginMessage = event.NewEventSpec[LoginEvent](
		AuditExchange, LoginRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
	LogoutMessage = event.NewEventSpec[LogoutEvent](
		AuditExchange, LogoutRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
	LikeMessage = event.NewEventSpec[LikeEvent](
		InteractionExchange, LikeRoutingKey,
		event.WithProducer(producer), event.WithVersion("v1"))
)

// 通知域按事件拆队列：站内信扇出可能长时间占用 worker，
// 与验证码邮件/短信共用队列会造成队头阻塞。
func newNotificationSpec(queue, bindingKey string, concurrency, maxInFlight int) mqx.SubscriptionSpec {
	return mqx.SubscriptionSpec{
		Destination: NotificationExchange,
		BindingKey:  bindingKey,
		Group:       queue,
		Concurrency: concurrency,
		MaxInFlight: maxInFlight,
		Retry: mqx.RetryPolicy{
			MaxAttempts:  3,
			InitialDelay: 200 * time.Millisecond,
			MaxDelay:     2 * time.Second,
			MaxElapsed:   10 * time.Second,
			Jitter:       0.2,
		},
		DeadLetter: mqx.DeadLetterPolicy{
			Enabled:     true,
			Destination: NotificationDLQ,
			// 三个通知订阅共用一个死信 exchange，绑定必须收窄到本事件；
			// 留空会退化成 #，同一条死信会被复制进每一条 dlq 队列。
			RoutingKey:  bindingKey,
			IncludeMeta: true,
		},
		OnExhausted: mqx.ExhaustedDeadLetter,
	}
}

// NewEmailSubscription 返回邮件队列订阅规格。
func NewEmailSubscription() mqx.SubscriptionSpec {
	return newNotificationSpec(EmailQueue, EmailRoutingKey, 2, 10)
}

// NewSmsSubscription 返回短信队列订阅规格。
func NewSmsSubscription() mqx.SubscriptionSpec {
	return newNotificationSpec(SmsQueue, SmsRoutingKey, 2, 10)
}

// NewInboxSubscription 返回站内信队列订阅规格。
func NewInboxSubscription() mqx.SubscriptionSpec {
	return newNotificationSpec(InboxQueue, InboxRoutingKey, 4, 20)
}

// Subscriptions 返回本进程要启动的全部订阅规格。
func Subscriptions() []mqx.SubscriptionSpec {
	return []mqx.SubscriptionSpec{
		NewEmailSubscription(),
		NewSmsSubscription(),
		NewInboxSubscription(),
		NewAuditSubscription(),
		NewInteractionSubscription(),
	}
}

// NewInteractionSubscription 返回点赞队列订阅规格。
// 点赞以 Redis 为实时状态，落库失败进入死信，不静默丢弃：
// 丢失的明细会让「谁点过」无法从库中重建。
func NewInteractionSubscription() mqx.SubscriptionSpec {
	return mqx.SubscriptionSpec{
		Destination: InteractionExchange,
		BindingKey:  "interaction.#",
		Group:       InteractionQueue,
		Concurrency: 4,
		MaxInFlight: 20,
		Retry: mqx.RetryPolicy{
			MaxAttempts:  3,
			InitialDelay: 200 * time.Millisecond,
			MaxDelay:     2 * time.Second,
			MaxElapsed:   10 * time.Second,
			Jitter:       0.2,
		},
		DeadLetter: mqx.DeadLetterPolicy{
			Enabled:     true,
			Destination: InteractionDLQ,
			RoutingKey:  "interaction.#",
			IncludeMeta: true,
		},
		OnExhausted: mqx.ExhaustedDeadLetter,
	}
}

// NewAuditSubscription 返回审计域订阅规格。
// 登录/登出日志是观测数据：失败重试后显式丢弃，不阻塞登录流程。
// 并发设为 1：登出依赖同一设备的登录记录，必须按 login -> logout 顺序处理。
// 该顺序只在「同一队列只有一个活跃消费者」时成立：若 blog-rpc 多副本，
// 同组消费者会分摊消息，需改用单活消费者或按设备键分区。
func NewAuditSubscription() mqx.SubscriptionSpec {
	return mqx.SubscriptionSpec{
		Destination: AuditExchange,
		BindingKey:  "audit.#",
		Group:       AuditQueue,
		Concurrency: 1,
		MaxInFlight: 4,
		Retry: mqx.RetryPolicy{
			MaxAttempts:  3,
			InitialDelay: 100 * time.Millisecond,
			MaxDelay:     time.Second,
			MaxElapsed:   5 * time.Second,
			Jitter:       0.2,
		},
		OnExhausted: mqx.ExhaustedDiscard,
	}
}
