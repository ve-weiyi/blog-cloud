package weibo

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
)

func TestWeibo(t *testing.T) {
	conf := &oauth.OauthConfig{
		ClientId:     "xxx",
		ClientSecret: "xxx",
		RedirectUri:  "https://blog.veweiyi.cn/oauth/login/weibo",
	}

	auth := NewAuthWb(conf)
	// 获取第三方登录地址
	url := auth.GetAuthLoginUrl("state")
	log.Println("url:", url)

	// 获取token信息
	// tokenRes, err := auth.GetAccessToken("9cfec63e61f53c7769a7c236827d85de")
	// log.Println("tokenRes:", err, tokenRes)
	//  {"access_token":"2.00OYpWYGPTpttCaf929b916cL6FMXD","remind_in":"157679999","expires_in":157679999,"uid":"6007017078","isRealName":"true"}
	// 刷新token
	// refresh, err := auth.RefreshToken("AE6F4302DD7AFB52902F56150FC58D2A")
	// log.Println("refresh:", err, refresh)

	// 获取用户信息
	userInfo, err := auth.GetAuthUserInfo("1c0bde90de5b063d3cf95887ca771f12")
	log.Println("userInfo:", err, userInfo)

}
