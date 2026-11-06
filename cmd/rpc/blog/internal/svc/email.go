package svc

import (
	"fmt"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/vkit/adapter/emailx"
)

// NewEmailDeliver 创建邮件通道。
func NewEmailDeliver(c config.EmailConf) (emailx.IEmailDeliver, error) {
	if c.Port == 465 && !c.IsSsl {
		return nil, fmt.Errorf("smtp port 465 requires is_ssl=true")
	}

	return emailx.NewEmailDeliver(toEmailConfig(c)), nil
}

func toEmailConfig(c config.EmailConf) *emailx.EmailConfig {
	return &emailx.EmailConfig{
		Host:     c.Host,
		Port:     c.Port,
		Username: c.Username,
		Password: c.Password,
		Nickname: c.Nickname,
		SSL:      c.IsSsl,
		BCC:      c.BCC,
	}
}
