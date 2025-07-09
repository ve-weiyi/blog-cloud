package feishu

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

func TestFeishu(t *testing.T) {
	conf := &oauth.OauthConfig{
		ClientId:     "cli_a5082b89657c900c",
		ClientSecret: "QJiJDotAX6mNLbBIhEQIocpbWza6el4J",
		RedirectUri:  "https://ankersolix-professional-ci.anker.com/login",
	}

	auth := NewAuthFeishu(conf)
	// 获取第三方登录地址
	url := auth.GetAuthLoginUrl("state")
	log.Println("url:", url)

	// 获取用户信息
	userInfo, err := auth.GetAuthUserInfo("fwBl0AAD5DzG4by09L7yxCa9CxBCHHe2")
	if err != nil {
		return
	}
	log.Println("userInfo:", err, userInfo)
}
