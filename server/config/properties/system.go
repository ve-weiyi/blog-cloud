package properties

type System struct {
	Version       string `mapstructure:"version" json:"version" yaml:"version"`                      // 程序版本
	Env           string `mapstructure:"env" json:"env" yaml:"env"`                                  // 环境值
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`                               // 端口
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`    // 路由前缀
	UseRedis      bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	RuntimePath   string `mapstructure:"runtime-path" json:"runtime-path" yaml:"runtime-path"`       // 运行时目录
}
