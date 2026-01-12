package gitee

import (
	"log"
	"net/url"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
)

func TestGitee(t *testing.T) {
	conf := &oauth.OauthConfig{
		ClientId:     "3836732746457306df93147721a0a38686536845d87e8807c7f63feca0770206",
		ClientSecret: "f61a0edadc75b8b3e7a0850509add17b04eb52f7195c8b35c5617f952f9d4777",
		RedirectUri:  "http://127.0.0.1:9421/#/oauth/login/gitee",
	}

	auth := NewAuthGitee(conf)
	// 获取第三方登录地址
	encodedStr := auth.GetAuthLoginUrl("state")
	log.Println("authorize url:", encodedStr)

	// 解码
	decodedStr, err := url.QueryUnescape(encodedStr)
	log.Println("decodedStr:", decodedStr)

	// 获取用户信息
	userInfo, err := auth.GetAuthUserInfo("0e7f1ee69736fd6b1b8e39c816c3975a790438235852747e3ff929aa049cc0b6")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("userInfo:", err, userInfo)
}
