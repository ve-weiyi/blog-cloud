package feishu

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/initest"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/oauth"
)

func TestFeishu(t *testing.T) {
	initest.InitConfig()
	conf := &oauth.AuthConfig{
		ClientID:     global.CONFIG.Oauth.Feishu.ClientID,
		ClientSecret: global.CONFIG.Oauth.Feishu.ClientSecret,
		RedirectUri:  "https://veweiyi.cn/blog/oauth/login/feishu",
	}

	auth := NewAuthFeishu(conf)
	//获取第三方登录地址
	url := auth.GetRedirectUrl("state")
	log.Println("url:", url)

	//获取用户信息
	userInfo, err := auth.GetUserOpenInfo("5ebhf5a54b19408892c18e2042edf792")
	if err != nil {
		return
	}
	log.Println("userInfo:", err, userInfo)
}
