package gitee

import (
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

func TestGitee(t *testing.T) {
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

	conf := ms["gitee"]
	log.Println("conf:", conf)
	auth := NewAuthGitee(conf)
	// 获取第三方登录地址
	url := auth.GetAuthorizeUrl("state")
	log.Println("authorize url:", url)

	// 获取用户信息
	userInfo, err := auth.GetUserOpenInfo("0e7f1ee69736fd6b1b8e39c816c3975a790438235852747e3ff929aa049cc0b6")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("userInfo:", err, userInfo)
}
