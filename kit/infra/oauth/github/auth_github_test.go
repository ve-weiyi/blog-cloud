package github

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
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

func TestName(t *testing.T) {
	openid := "67481255"
	uid := crypto.Md5v(openid, "")
	t.Log("uid:", uid)
	t.Log(base64.StdEncoding.EncodeToString([]byte(uid)))
}
