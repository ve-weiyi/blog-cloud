package github

import (
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

func TestGithub(t *testing.T) {
	v := viper.New()
	v.SetConfigFile("../config.yaml")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	ms := map[string]*oauth.AuthConfig{}
	err = v.Unmarshal(&ms, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	})
	if err != nil {
		log.Fatal(err)
	}

	conf := ms["github"]
	auth := NewAuthGithub(conf)
	// 获取第三方登录地址
	url := auth.GetAuthorizeUrl("state")
	log.Println("url:", url)

	// 获取用户信息
	userInfo, err := auth.GetUserOpenInfo("69b1ccc128ef6aada3c4")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("userInfo:", err, userInfo)
}
