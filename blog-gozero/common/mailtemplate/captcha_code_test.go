package mailtemplate

import (
	_ "embed"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/quickstart/gotplgen"
)

func Test_Mail(t *testing.T) {
	data := CaptchaEmail{
		Username: "791422171@qq.com",
		Content:  "欢迎注册我的博客平台。",
		Code:     "123456",
	}

	meta := gotplgen.TemplateMeta{
		Mode:           gotplgen.ModeCreateOrReplace,
		CodeOutPath:    "./runtime/code.html",
		TemplateString: TempCaptchaCode,
		FunMap:         nil,
		Data:           data,
	}

	err := meta.Execute()
	t.Log(err)
}
