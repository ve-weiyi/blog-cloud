package mailtemplate

import (
	_ "embed"
	"testing"

	"github.com/ve-weiyi/vkit/gen/tmplx"
)

func Test_Mail(t *testing.T) {
	data := CaptchaEmail{
		Username: "791422171@qq.com",
		Content:  "欢迎注册我的博客平台。",
		Code:     "123456",
	}

	meta := tmplx.TemplateMeta{
		Mode:           tmplx.ModeCreateOrReplace,
		CodeOutPath:    "./runtime/code.html",
		TemplateString: TempCaptchaCode,
		FunMap:         nil,
		Data:           data,
	}

	err := meta.Execute()
	t.Log(err)
}
