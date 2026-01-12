package mail

type Option func(*EmailConfig)

func WithHost(host string) Option {
	return func(o *EmailConfig) {
		o.Host = host
	}
}

func WithPort(port int) Option {
	return func(o *EmailConfig) {
		o.Port = port
	}
}

func WithUsername(username string) Option {
	return func(o *EmailConfig) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *EmailConfig) {
		o.Password = password
	}
}

func WithNickname(nickname string) Option {
	return func(o *EmailConfig) {
		o.Nickname = nickname
	}
}

func WithDeliver(deliver []string) Option {
	return func(o *EmailConfig) {
		o.CC = deliver
	}
}

func WithSSL(ssl bool) Option {
	return func(o *EmailConfig) {
		o.SSL = ssl
	}
}
