package mail

import (
	"crypto/tls"
	"log"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
	"gopkg.in/gomail.v2"
)

func Test_Send(t *testing.T) {

	host := "smtp.qq.com"
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "647166282@qq.com", "与梦")
	m.SetHeader("To", "791422171@qq.com")
	m.SetHeader("Cc", "919390162@qq.com")
	m.SetHeader("Bcc", "791422171@qq.com")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")

	d := gomail.NewDialer(host, 465, "647166282@qq.com", "culyqmzvmppabccd")

	s, err := d.Dial()
	if err != nil {
		panic(err)
	}
	// Send the email to Bob, Cora and Dan.
	if err := gomail.Send(s, m); err != nil {
		panic(err)
	}
}

func Test_Send2(t *testing.T) {
	host := "smtp.qq.com"

	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "与梦 <647166282@qq.com>"
	// 设置接收方的邮箱
	e.To = []string{"791422171@qq.com"}
	//设置抄送如果抄送多人逗号隔开
	e.Cc = []string{"919390162@qq.com"}
	//设置秘密抄送
	e.Bcc = []string{"919390162@qq.com"}
	//设置主题
	e.Subject = "这是主题"
	//设置文件发送的内容
	e.Text = []byte("这是正文")

	auth := smtp.PlainAuth("", "647166282@qq.com", "culyqmzvmppabccd", host)
	//设置服务器相关的配置
	err := e.SendWithTLS("smtp.qq.com:465", auth, &tls.Config{ServerName: "smtp.qq.com"})
	if err != nil {
		log.Fatal(err)
	}
}
