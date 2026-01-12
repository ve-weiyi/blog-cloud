package github

import (
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
)

func TestGithub(t *testing.T) {
	conf := &oauth.OauthConfig{
		ClientId:     "Ov23li2CTGk4hHq93ZYz",
		ClientSecret: "73804611edc3f1f86ead487a189dd43b1fffaf76",
		RedirectUri:  "http://127.0.0.1:9421/oauth/login/github",
	}

	auth := NewAuthGithub(conf)
	// 获取第三方登录地址
	url := auth.GetAuthLoginUrl("")
	log.Println("url:", url)

	// 获取用户信息
	userInfo, err := auth.GetAuthUserInfo("97b175340da66743a50d")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("userInfo:", err, userInfo)
}
