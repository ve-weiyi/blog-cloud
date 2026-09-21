package mq

// ── 交换机 ──
// 每种消息一个 fanout 交换机，交换机名与队列名在 broker 上必须保持一致，
// 重命名会导致旧交换机/队列成为孤儿。

const (
	EmailExchange  = "blog-email-exchange"
	SmsExchange    = "blog-sms-exchange"
	InboxExchange  = "blog-inbox-exchange"
	LoginExchange  = "blog-login-log-exchange"
	LogoutExchange = "blog-logout-exchange"
)

// ── 路由键 ──
// fanout 交换机忽略路由键，此处的值仅作为队列绑定键保留。

const (
	EmailRoutingKey  = "email"
	SmsRoutingKey    = "sms"
	InboxRoutingKey  = "inbox"
	LoginRoutingKey  = "login"
	LogoutRoutingKey = "logout"
)

// ── 队列名（消费端私有）──

const (
	EmailQueue  = "blog-email-queue"
	SmsQueue    = "blog-sms-queue"
	InboxQueue  = "blog-inbox-queue"
	LoginQueue  = "blog-login-log-queue"
	LogoutQueue = "blog-logout-queue"
)

// ── 事件 ──

// LoginEvent 登录日志事件
type LoginEvent struct {
	UserId     string `json:"user_id"`
	DeviceId   string `json:"device_id"`
	LoginType  string `json:"login_type"`
	Status     int64  `json:"status"`
	FailReason string `json:"fail_reason,omitempty"`
	// 游客信息
	IpAddress string `json:"ip_address,omitempty"`
	Os        string `json:"os,omitempty"`
	Browser   string `json:"browser,omitempty"`
	Device    string `json:"device,omitempty"`
	Location  string `json:"location,omitempty"`
}

// LogoutEvent 登出事件
type LogoutEvent struct {
	UserId     string `json:"user_id"`
	DeviceId   string `json:"device_id"`
	LogoutType string `json:"logout_type"` // manual-手动登出 timeout-超时登出 force-强制登出
}

// EmailMessageEvent 邮件消息事件
type EmailMessageEvent struct {
	Email   string            `json:"email"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Scene   string            `json:"scene"`
	BizId   string            `json:"biz_id"`
	Params  map[string]string `json:"params,omitempty"`
}

// SmsMessageEvent SMS 消息事件
type SmsMessageEvent struct {
	Mobile string            `json:"mobile"`
	Scene  string            `json:"scene"`
	BizId  string            `json:"biz_id"`
	Params map[string]string `json:"params,omitempty"`
}

// InboxMessageEvent 站内信投递事件
// 消息发布后，异步为每个目标用户创建 delivery 记录
type InboxMessageEvent struct {
	MessageId int64 `json:"message_id"`
}
