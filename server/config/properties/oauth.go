package properties

type Oauth struct {
	QQ     AuthConfig `json:"qq" yaml:"qq" mapstructure:"qq"`
	Feishu AuthConfig `json:"feishu" yaml:"feishu" mapstructure:"feishu"`
	Weibo  AuthConfig `json:"weibo" yaml:"weibo" mapstructure:"weibo"`
}

type AuthConfig struct {
	ClientID     string `json:"client-id" yaml:"client-id" mapstructure:"client-id"`
	ClientSecret string `json:"client-secret" yaml:"client-secret" mapstructure:"client-secret"`
	RedirectUri  string `json:"redirect-uri" yaml:"redirect-uri" mapstructure:"redirect-uri"`
}
