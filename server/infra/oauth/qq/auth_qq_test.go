package qq

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/initest"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/oauth"
)

func TestQQ(t *testing.T) {
	initest.InitConfig()
	conf := &oauth.AuthConfig{
		ClientID:     global.CONFIG.Oauth.QQ.ClientID,
		ClientSecret: global.CONFIG.Oauth.QQ.ClientSecret,
		RedirectUri:  global.CONFIG.Oauth.QQ.RedirectUri,
	}

	auth := NewAuthQq(conf)
	//获取第三方登录地址
	url := auth.GetRedirectUrl("state")
	log.Println("url:", url)
	//获取用户信息
	//userInfo, err := auth.GetUserOpenInfo("D3337DCCFF3A9ACD1A3F4501E90AC7F5")
	//log.Println("userInfo:", err, userInfo)
	//
	////获取token信息
	//tokenRes, err := auth.GetAccessToken("D3337DCCFF3A9ACD1A3F4501E90AC7F5")
	//log.Println("tokenRes:", err, tokenRes)

	//刷新token
	refresh, err := auth.RefreshToken("9AEA28A71B91AF087CB6D3986BA62D24")
	log.Println("refresh:", err, refresh)
}
