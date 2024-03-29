package config

import (
	"github.com/ve-weiyi/ve-blog-golang/server/config/properties"
)

type Config struct {
	System properties.System `mapstructure:"system" json:"system" yaml:"system"`
	Mysql  properties.Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  properties.Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	// 消息队列、邮件发送服务
	RabbitMQ properties.RabbitMQ `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
	Email    properties.Email    `mapstructure:"email" json:"email" yaml:"email"`
	// 鉴权
	JWT   properties.JWT   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Oauth properties.Oauth `mapstructure:"oauth" json:"oauth" yaml:"oauth"`
	// 日志
	Zap properties.Zap `mapstructure:"zap" json:"zap" yaml:"zap"`
	// 文件上传
	Upload  properties.Upload  `mapstructure:"upload" json:"upload" yaml:"upload"`
	Captcha properties.Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	ChatGPT properties.ChatGPT `mapstructure:"chatgpt" json:"chatgpt" yaml:"chatgpt"`
}
