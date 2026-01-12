package svc

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/mail"
)

func Test_Captcha(t *testing.T) {
	sv := NewTestServiceContext()
	key := "test"
	captcha, err := sv.CaptchaHolder.GetCodeCaptcha(key)
	if err != nil {
		log.Fatal(err)
	}

	t.Log(captcha)
	if sv.CaptchaHolder.VerifyCaptcha(key, captcha) {
		t.Log("verify success")
	}

}

func TestInitEmailDeliver(t *testing.T) {
	sv := NewTestServiceContext()

	msg := &mail.EmailMessage{
		To:      []string{"791422171@qq.com"},
		Subject: "测试邮件标题",
		Content: "测试邮件内容",
	}

	err := sv.EmailDeliver.DeliveryEmail(msg)
	if err != nil {
		t.Errorf("send email error: %v", err)
	}

	t.Log("send email success")

	select {}
}
