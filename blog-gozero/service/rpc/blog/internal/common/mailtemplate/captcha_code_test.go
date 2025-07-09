package mailtemplate

import (
	_ "embed"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/quickstart/invent"
)

func Test_Mail(t *testing.T) {
	data := CaptchaEmail{
		Username: "791422171@qq.com",
		Content:  "欢迎注册我的博客平台。",
		Code:     "123456",
	}

	meta := invent.TemplateMeta{
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    "./runtime/code.html",
		TemplateString: TempCaptchaCode,
		FunMap:         nil,
		Data:           data,
	}

	err := meta.Execute()
	t.Log(err)
}
