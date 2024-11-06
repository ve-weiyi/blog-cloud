package oauth

type Option func(*AuthConfig)

func WithClientID(clientId string) Option {
	return func(o *AuthConfig) {
		o.ClientId = clientId
	}
}

func WithClientSecret(clientSecret string) Option {
	return func(o *AuthConfig) {
		o.ClientSecret = clientSecret
	}
}

func WithRedirectUri(redirectUri string) Option {
	return func(o *AuthConfig) {
		o.RedirectUri = redirectUri
	}
}
