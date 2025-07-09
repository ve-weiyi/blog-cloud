package mailtemplate

import (
	_ "embed"
)

// 验证码邮件内容
type CaptchaEmail struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Code     string `json:"code"`
}

//go:embed captcha_code.tpl
var TempCaptchaCode string
