package config

type Config struct {
	System   SystemConf   `json:"system" yaml:"system"`
	JWT      JWTConf      `json:"jwt" yaml:"jwt"`
	Zap      Zap          `json:"zap" yaml:"zap"`
	Mysql    MysqlConf    `json:"mysql" yaml:"mysql"`
	Redis    RedisConf    `json:"redis" yaml:"redis"`
	RabbitMQ RabbitMQConf `json:"rabbitmq" yaml:"rabbitmq"`
}

// 系统配置
type SystemConf struct {
	Version     string `json:"version" yaml:"version"`           // 程序版本
	Env         string `json:"env" yaml:"env"`                   // 运行环境
	Port        int    `json:"port" yaml:"port"`                 // 运行端口
	RuntimePath string `json:"runtime-path" yaml:"runtime-path"` // 运行时目录
}

// jwt鉴权
type JWTConf struct {
	SigningKey  string `json:"signing-key" yaml:"signing-key"`   // jwt签名
	ExpiresTime string `json:"expires-time" yaml:"expires-time"` // 过期时间
	Issuer      string `json:"issuer" yaml:"issuer"`             // 签发者
	Type        string `json:"type" yaml:"type"`                 // 类型 例如：Bearer
}

// mysql数据库
type MysqlConf struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Dbname   string `json:"dbname" yaml:"dbname"`
	Config   string `json:"config" yaml:"config"`
}

// redis缓存
type RedisConf struct {
	Host     string `json:"host" yaml:"host"` // 服务器地址:端口
	Port     string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"` // 密码
	DB       int    `json:"db" yaml:"db"`             // redis的哪个数据库
}

// rabbitmq消息队列
type RabbitMQConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Zap struct {
	Mode           string `json:"mode" yaml:"mode"`                       // 模式
	Format         string `json:"format" yaml:"format"`                   // 输出
	Level          string `json:"level" yaml:"level"`                     // 级别
	Prefix         string `json:"prefix" yaml:"prefix"`                   // 日志前缀
	EncodeCaller   string `json:"encode-caller" yaml:"encode-caller"`     // 编码调用者
	EncodeColorful bool   `json:"encode-colorful" yaml:"encode-colorful"` // 编码调用者

	Filename string `json:"filename" yaml:"filename"`   // 日志文件名称
	CacheDir string `json:"cache-dir" yaml:"cache-dir"` // 日志文件夹
	MaxAge   int    `json:"max-age" yaml:"max-age"`     // 日志留存时间
}
