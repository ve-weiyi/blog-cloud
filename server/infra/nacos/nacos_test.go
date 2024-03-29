package nacos

import (
	"fmt"
	"testing"
)

type ConfigReloaderMock struct {
}

func (sel *ConfigReloaderMock) Reload(data string) error {
	fmt.Println(data)
	return nil
}

func TestNewNacos(t *testing.T) {

	cfg := &NacosConfig{
		IP:          "120.79.136.81",
		Port:        8848,
		DataID:      "ve-blog-golang",
		Group:       "veweiyi.cn",
		UserName:    "nacos",
		Password:    "nacos",
		LogLevel:    "warn",
		NameSpaceID: "blog",
		Timeout:     5000,
	}

	err := New(cfg).Init(ConfigReloaderMock{}.Reload)
	fmt.Println(err)
}
