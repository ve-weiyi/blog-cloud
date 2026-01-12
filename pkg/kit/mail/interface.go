package mail

type EmailConfig struct {
	Host     string   `json:"host"`     // 服务器地址
	Port     int      `json:"port"`     // 端口
	Username string   `json:"username"` // 发件人
	Password string   `json:"password"` // 密钥
	Nickname string   `json:"nickname"` // 发件人昵称
	SSL      bool     `json:"ssl"`      // 是否加密
	CC       []string `json:"cc"`       // 抄送邮箱:多个以英文逗号分隔
	BCC      []string `json:"bcc"`      // 密送邮箱:多个以英文逗号分隔
}

type EmailMessage struct {
	To      []string `json:"to"`      // 目标邮箱号
	CC      []string `json:"cc"`      // 抄送人
	Subject string   `json:"subject"` // 主题
	Content string   `json:"content"` // 内容
}

// 投递邮件
type IEmailDeliver interface {
	DeliveryEmail(message *EmailMessage) error
}
