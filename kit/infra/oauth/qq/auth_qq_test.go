package qq

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

func TestQQ(t *testing.T) {
	conf := &oauth.AuthConfig{
		ClientId:     "xxx",
		ClientSecret: "xxx",
		RedirectUri:  "https://veweiyi.cn/blog/oauth/login/qq",
	}

	auth := NewAuthQq(conf)
	// 获取第三方登录地址
	url := auth.GetAuthorizeUrl("state")
	log.Println("url:", url)
	// 获取用户信息
	// userInfo, err := auth.GetUserOpenInfo("D3337DCCFF3A9ACD1A3F4501E90AC7F5")
	// log.Println("userInfo:", err, userInfo)
	//
	// //获取token信息
	// tokenRes, err := auth.GetAccessToken("D3337DCCFF3A9ACD1A3F4501E90AC7F5")
	// log.Println("tokenRes:", err, tokenRes)

	// 刷新token
	refresh, err := auth.RefreshToken("9AEA28A71B91AF087CB6D3986BA62D24")
	log.Println("refresh:", err, refresh)
}
