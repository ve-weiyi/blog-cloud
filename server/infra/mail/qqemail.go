package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender struct {
	Host     string   // 服务器地址
	Port     int      // 端口
	Username string   // 发件人
	Password string   // 密钥
	Nickname string   // 发件人昵称
	Deliver  []string // 抄送邮箱:多个以英文逗号分隔
	IsSSL    bool     // 是否使用 SSL/TLS
}

func (s *EmailSender) SendEmailMessage(message EmailMessage) error {
	return s.send(message.To, message.Subject, message.Content, message.Type == 1)
}

// 发送邮件
// ReplyTo: 回复邮件时，邮件的回复地址。可以是多个地址。
// From: 发件人的邮箱地址。
// BindExchange: 收件人的邮箱地址。可以是多个地址。
// Bcc: 密送（暗送）的邮箱地址。收件人不可见。可以是多个地址。
// Cc: 抄送的邮箱地址。收件人可见。可以是多个地址。
// Subject: 邮件的主题。
// Text: 邮件的纯文本内容，可选字段。
// HTML: 邮件的 HTML 内容，可选字段。
// Sender: 覆盖发件人的 SMTP 信封发送者，可选字段。
// Headers: 邮件的附加头部信息，使用 textproto.MIMEHeader 类型存储。
// Attachments: 邮件的附件列表，每个附件是一个 Attachment 结构体的实例。
// ReadReceipt: 邮件的回执邮箱地址，表示邮件阅读回执的接收地址。可以是多个地址。
func (s *EmailSender) send(to []string, subject string, body string, cc bool) (err error) {
	host := s.Host
	port := s.Port
	ssl := s.IsSSL
	username := s.Username
	password := s.Password
	nickname := s.Nickname
	deliver := s.Deliver

	auth := smtp.PlainAuth("", username, password, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, username)
	} else {
		e.From = username
	}
	e.To = to
	// 抄送
	if cc {
		e.Cc = deliver
	}
	e.Subject = subject
	e.HTML = []byte(body)
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if ssl {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
