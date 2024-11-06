/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/nacos"
	"github.com/ve-weiyi/ve-blog-golang/server/config"
	"github.com/ve-weiyi/ve-blog-golang/server/core"
)

type ApiCmd struct {
	cmd        *cobra.Command
	configMode string // 运行方式 file|nacos
	filepath   string
	nacosCfg   *nacos.NacosConfig
}

func NewApiCmd() *ApiCmd {
	apiCmd := &ApiCmd{}
	apiCmd.cmd = &cobra.Command{
		Use:   "api",
		Short: "启动接口服务",
		Long:  `启动接口服务`,
		Run: func(cmd *cobra.Command, args []string) {
			apiCmd.RunApi(cmd, args)
		},
	}

	apiCmd.init()
	return apiCmd
}

func (s *ApiCmd) init() {
	nacosCfg := s.GetDefaultNacosConfig()
	s.nacosCfg = nacosCfg
	// 设置默认参数
	s.cmd.Flags().StringVarP(&s.configMode, "config", "c", "file", "the way of read config file (file|nacos)")
	s.cmd.Flags().StringVarP(&s.filepath, "filepath", "f", "config.yaml", "config file path (default is ./config.yaml)")
	s.cmd.Flags().StringVar(&s.nacosCfg.IP, "n-ip", nacosCfg.IP, "the ip for nacos")
	s.cmd.Flags().Uint64Var(&s.nacosCfg.Port, "n-port", nacosCfg.Port, "the port for nacos")
	s.cmd.Flags().StringVar(&s.nacosCfg.UserName, "n-user", nacosCfg.UserName, "the user for nacos")
	s.cmd.Flags().StringVar(&s.nacosCfg.Password, "n-password", nacosCfg.Password, "the password for nacos")
	s.cmd.Flags().StringVar(&s.nacosCfg.DataId, "n-data-id", nacosCfg.DataId, "the DataId for nacos")
	s.cmd.Flags().StringVar(&s.nacosCfg.Group, "n-group", nacosCfg.Group, "the group for nacos")
	s.cmd.Flags().StringVar(&s.nacosCfg.NameSpaceId, "n-namespace", nacosCfg.NameSpaceId, "the namespace for nacos")
}

func (s *ApiCmd) GetDefaultNacosConfig() *nacos.NacosConfig {
	return &nacos.NacosConfig{
		IP:          "120.79.136.81",
		Port:        8848,
		UserName:    "nacos",
		Password:    "nacos",
		NameSpaceId: "dev",
		Group:       "veweiyi.cn",
		DataId:      "ve-blog-golang",
		RuntimeDir:  "runtime/log/nacos",
		LogLevel:    "warn",
		Timeout:     5000,
	}
}

func (s *ApiCmd) RunApi(cmd *cobra.Command, args []string) {
	var c config.Config
	var content string

	switch s.configMode {
	case "file":
		log.Println("读取配置文件..使用文件路径")

		text, err := os.ReadFile(s.filepath)
		if err != nil {
			panic(err)
		}

		content = string(text)

	case "nacos":
		log.Println("读取配置文件...使用nacos")
		// 初始化Nacos
		nc := nacos.New(s.nacosCfg)

		// 读取配置文件
		text, err := nc.GetConfig()
		if err != nil {
			panic("nacos config read failed " + err.Error())
		}

		content = text
	default:
		panic("config mode not support,please use cmd 'go run main.go api --c=file --f=./config.yaml'")
	}

	// 初始化Viper
	v := viper.New()
	v.SetConfigType("yaml")

	// 读取配置文件
	err := v.ReadConfig(strings.NewReader(content))
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	// 修改解析的tag（默认是mapstructure）
	withJsonTag := func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"
	}
	// 解析配置文件
	if err = v.Unmarshal(&c, withJsonTag); err != nil {
		panic(err)
	}
	// 暂时不开启监听配置文件变化

	// 初始化配置文件
	core.RunWindowsServer(&c)
}
