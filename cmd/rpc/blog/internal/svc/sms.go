package svc

import (
	"fmt"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/vkit/adapter/smsx"
)

// NewSmsProvider 创建短信通道；非 mock 提供方必须配齐各场景模板。
func NewSmsProvider(c config.SmsConf) (smsx.SmsProvider, error) {
	if c.Provider != "mock" {
		for _, scene := range []string{enums.CodeSceneMobileLogin, enums.CodeSceneResetPassword, enums.CodeSceneBind} {
			if c.Templates[scene] == "" {
				return nil, fmt.Errorf("sms template missing for scene %s", scene)
			}
		}
	}

	return smsx.NewSmsProvider(&smsx.SmsConfig{
		Provider:  c.Provider,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
		SignName:  c.SignName,
		Region:    c.Region,
		SdkAppId:  c.SdkAppId,
		Templates: c.Templates,
	})
}
