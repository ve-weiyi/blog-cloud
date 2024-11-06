package feishu

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

func TestFeishu(t *testing.T) {
	conf := &oauth.AuthConfig{
		ClientId:     "xxx",
		ClientSecret: "xxx",
		RedirectUri:  "https://veweiyi.cn/blog/oauth/login/feishu",
	}

	auth := NewAuthFeishu(conf)
	// 获取第三方登录地址
	url := auth.GetAuthorizeUrl("state")
	log.Println("url:", url)

	// 获取用户信息
	userInfo, err := auth.GetUserOpenInfo("5ebhf5a54b19408892c18e2042edf792")
	if err != nil {
		return
	}
	log.Println("userInfo:", err, userInfo)
}
