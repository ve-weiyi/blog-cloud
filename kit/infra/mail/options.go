package mail

type Option func(*EmailDeliver)

func WithHost(host string) Option {
	return func(o *EmailDeliver) {
		o.Host = host
	}
}

func WithPort(port int) Option {
	return func(o *EmailDeliver) {
		o.Port = port
	}
}

func WithUsername(username string) Option {
	return func(o *EmailDeliver) {
		o.Username = username
	}
}

func WithPassword(password string) Option {
	return func(o *EmailDeliver) {
		o.Password = password
	}
}

func WithNickname(nickname string) Option {
	return func(o *EmailDeliver) {
		o.Nickname = nickname
	}
}

func WithDeliver(deliver []string) Option {
	return func(o *EmailDeliver) {
		o.Deliver = deliver
	}
}

func WithSSL(ssl bool) Option {
	return func(o *EmailDeliver) {
		o.IsSSL = ssl
	}
}
