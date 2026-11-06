package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	MysqlConf    MysqlConf
	RedisConf    RedisConf
	RabbitMQConf RabbitMQConf
	EmailConf    EmailConf
	SmsConf      SmsConf

	AppOAuthConf map[string]OAuthConf // map[platform]conf
	NotifyConf   NotifyConf
}

// 通知配置
type NotifyConf struct {
	// UnsubscribeSecret 邮件退订令牌的签名密钥。为空时退订接口直接拒绝，
	// 不给「没有签名也能退订」留后门。
	UnsubscribeSecret string `json:"unsubscribe_secret,optional"`

	// SiteURL 博客站点根地址，用于拼接邮件里的退订链接；为空时不渲染链接
	SiteURL string `json:"site_url,optional"`
}

// mysql数据库配置
type MysqlConf struct {
	Host     string `json:"host"`
	Port     string `json:"port,default=3306"`
	Username string `json:"username"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Config   string `json:"config,optional"` // 默认值包含特殊字符，在代码中处理
}

// redis缓存配置
type RedisConf struct {
	DB       int    `json:"db,default=0"`      // redis的哪个数据库
	Host     string `json:"host"`              // 服务器地址
	Port     string `json:"port,default=6379"` // 端口
	Password string `json:"password,optional"` // 密码
}

// rabbitmq配置
type RabbitMQConf struct {
	Host     string `json:"host"`
	Port     string `json:"port,default=5672"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 邮件配置
type EmailConf struct {
	Host     string   `json:"host"`                // 服务器地址
	Port     int      `json:"port,default=465"`    // 端口
	Username string   `json:"username"`            // 发件人
	Password string   `json:"password"`            // 密钥
	Nickname string   `json:"nickname,optional"`   // 发件人昵称
	BCC      []string `json:"bcc,optional"`        // 密送邮箱
	IsSsl    bool     `json:"is_ssl,default=true"` // 是否启用SSL
}

// 短信配置
type SmsConf struct {
	Provider  string            `json:"provider,default=mock"` // 服务商类型：aliyun | tencent | mock
	AccessKey string            `json:"access_key,optional"`   // AccessKey / SecretId
	SecretKey string            `json:"secret_key,optional"`   // SecretKey
	SignName  string            `json:"sign_name,optional"`    // 短信签名
	Region    string            `json:"region,optional"`       // 地域（腾讯云需要）
	SdkAppId  string            `json:"sdk_app_id,optional"`   // SDK应用ID（腾讯云需要）
	Templates map[string]string `json:"templates,optional"`    // 模板配置：codeType -> templateCode/templateId
}

type OAuthConf struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri,optional"`
}
