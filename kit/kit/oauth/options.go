package oauth

type Option func(*OauthConfig)

func WithClientID(clientId string) Option {
	return func(o *OauthConfig) {
		o.ClientId = clientId
	}
}

func WithClientSecret(clientSecret string) Option {
	return func(o *OauthConfig) {
		o.ClientSecret = clientSecret
	}
}

func WithRedirectUri(redirectUri string) Option {
	return func(o *OauthConfig) {
		o.RedirectUri = redirectUri
	}
}
