package mail

type EmailMessage struct {
	To      []string `json:"to"`      // 目标邮箱号
	Subject string   `json:"subject"` // 主题
	Content string   `json:"content"` // 内容
	CC      bool     `json:"cc"`      // 是否抄送
}

// 投递邮件
type IEmailDeliver interface {
	DeliveryEmail(message *EmailMessage) error
}
