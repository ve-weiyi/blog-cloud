package oauth

// 基本配置
type OauthConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

// 获取第三方登录地址
type Oauth interface {
	GetName() string
	GetAuthLoginUrl(state string) string
	GetAuthUserInfo(code string) (*UserResult, error)
	// GetAccessToken(code string) (*result.TokenResult, error)
	// RefreshToken(refreshToken string) (*result.RefreshResult, error)
}
