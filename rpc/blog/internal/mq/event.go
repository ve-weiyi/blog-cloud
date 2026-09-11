package mq

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
