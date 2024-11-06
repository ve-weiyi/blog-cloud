package chatgpt

type Config struct {
	ApiHost string
	ApiKey  string
	Model   string
}

type Option func(*Config)

func WithApiHost(apiHost string) Option {
	return func(o *Config) {
		o.ApiHost = apiHost
	}
}

func WithApiKey(apiKey string) Option {
	return func(o *Config) {
		o.ApiKey = apiKey
	}
}

func WithModel(model string) Option {
	return func(o *Config) {
		o.Model = model
	}
}
