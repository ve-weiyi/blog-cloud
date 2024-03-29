package properties

type ChatGPT struct {
	ApiHost string `mapstructure:"api-host" json:"api-host" yaml:"api-host"` // host
	ApiKey  string `mapstructure:"api-key" json:"api-key" yaml:"api-key"`    // 秘钥
	Model   string `mapstructure:"model" json:"model" yaml:"model"`          // 模型
}
