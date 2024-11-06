package oauth

// 基本配置
type AuthConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
}

type AuthOauthURL struct {
	// Config *AuthConfig //配置信息
	Name            string // 第三方名称
	AuthorizeUrl    string // 授权登录URL
	AccessTokenUrl  string // 获得访问令牌URL
	RefreshTokenUrl string // 刷新令牌URL
	OpenidUrl       string // 获取用户OPENID
	UserInfoUrl     string // 获取用户信息URL
}

// 获取第三方登录地址
type Oauth interface {
	GetName() string
	GetAuthorizeUrl(state string) string
	GetUserOpenInfo(code string) (*UserResult, error)
	// GetAccessToken(code string) (*result.TokenResult, error)
	// RefreshToken(refreshToken string) (*result.RefreshResult, error)
}
