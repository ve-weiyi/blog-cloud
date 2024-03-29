/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ve-weiyi/ve-blog-golang/server/initialize"
)

type RabbitmqCmd struct {
	cmd        *cobra.Command
	configFile string
}

func NewRabbitmqCmd() *RabbitmqCmd {
	cmdRabbitmq := &RabbitmqCmd{}
	cmdRabbitmq.cmd = &cobra.Command{
		Use:   "rabbitmq",
		Short: "运行rabbitmq服务",
		Long:  `运行rabbitmq服务，订阅消息,发送邮件`,
		Run: func(cmd *cobra.Command, args []string) {
			initialize.RabbitMq()
		},
	}
	cmdRabbitmq.init()
	return cmdRabbitmq
}

func (s *RabbitmqCmd) init() {
	// 设置默认参数
}
