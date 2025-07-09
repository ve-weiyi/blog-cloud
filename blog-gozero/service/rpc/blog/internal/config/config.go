package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MysqlConf      MysqlConf
	RedisConf      RedisConf
	RabbitMQConf   RabbitMQConf
	EmailConf      EmailConf
	ThirdPartyConf map[string]map[string]ThirdPartyInfo

	DefaultCommentReviewStatus int64 `json:",default=1"` // 默认评论审核状态
	DefaultRemarkReviewStatus  int64 `json:",default=1"` // 默认留言审核状态
}

// mysql数据库配置
type MysqlConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Config   string `json:"config"`
}

// redis缓存配置
type RedisConf struct {
	DB       int    `json:"db" yaml:"db"`     // redis的哪个数据库
	Host     string `json:"host" yaml:"host"` // 服务器地址:端口
	Port     string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"` // 密码
}

// rabbitmq配置
type RabbitMQConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 邮件配置
type EmailConf struct {
	Host     string   `json:"host"`     // 服务器地址
	Port     int      `json:"port"`     // 端口
	Username string   `json:"username"` // 发件人
	Password string   `json:"password"` // 密钥
	Nickname string   `json:"nickname"` // 发件人昵称
	BCC      []string `json:"bcc"`      // 密送邮箱:多个以英文逗号分隔
}

type ThirdPartyInfo struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}
