package oauth

// 基本配置
type AuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

type AuthOauth struct {
	Config         *AuthConfig //配置信息
	AuthorizeUrl   string      //授权登录URL
	TokenUrl       string      //获得令牌URL
	AccessTokenUrl string      //获得访问令牌URL
	RefreshUrl     string      //刷新令牌URL
	OpenidUrl      string      //获取用户OPENID
	UserInfoUrl    string      //获取用户信息URL
	registerSource string      //注册来源
}

func (b *AuthOauth) Set(sourceName RegisterSource, cfg *AuthConfig) {
	b.Config = cfg
	b.registerSource = string(sourceName)
}

// 获取第三方登录地址
type Oauth interface {
	GetRedirectUrl(state string) string
	GetUserOpenInfo(code string) (*UserResult, error)
	//GetAccessToken(code string) (*result.TokenResult, error)
	//RefreshToken(refreshToken string) (*result.RefreshResult, error)
}
