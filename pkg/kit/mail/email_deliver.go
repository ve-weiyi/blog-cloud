package mail

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type EmailDeliver struct {
	EmailConfig
}

func NewEmailDeliver(conf *EmailConfig, opts ...Option) *EmailDeliver {
	for _, opt := range opts {
		opt(conf)
	}

	return &EmailDeliver{
		EmailConfig: *conf,
	}
}

func (s *EmailDeliver) DeliveryEmail(message *EmailMessage) error {
	return s.send(message.To, message.CC, message.Subject, message.Content)
}

// 发送邮件
// ReplyTo: 回复邮件时，邮件的回复地址。可以是多个地址。
// From: 发件人的邮箱地址。
// To: 收件人的邮箱地址。可以是多个地址。
// Bcc: 密送（暗送）的邮箱地址。收件人不可见。可以是多个地址。
// Cc: 抄送的邮箱地址。收件人可见。可以是多个地址。
// Subject: 邮件的主题。
// Text: 邮件的纯文本内容，可选字段。
// HTML: 邮件的 HTML 内容，可选字段。
// Sender: 覆盖发件人的 SMTP 信封发送者，可选字段。
// Headers: 邮件的附加头部信息，使用 textproto.MIMEHeader 类型存储。
// Attachments: 邮件的附件列表，每个附件是一个 Attachment 结构体的实例。
// ReadReceipt: 邮件的回执邮箱地址，表示邮件阅读回执的接收地址。可以是多个地址。
func (s *EmailDeliver) send(to []string, cc []string, subject string, body string) (err error) {
	host := s.Host
	port := s.Port
	username := s.Username
	password := s.Password
	nickname := s.Nickname
	ssl := s.SSL

	m := gomail.NewMessage()
	m.SetAddressHeader("From", username, nickname)
	m.SetHeader("To", to...)
	m.SetHeader("ReplyTo", username)
	if len(cc) != 0 {
		m.SetHeader("Cc", cc...)
	}
	if len(s.BCC) != 0 {
		m.SetHeader("Bcc", s.BCC...)
	}
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)
	if ssl {
		d.SSL = true
		d.TLSConfig = &tls.Config{}
	}
	return d.DialAndSend(m)
}
