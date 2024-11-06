package qq

import (
	"encoding/json"
	"log"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/httpx"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
)

// QQ授权登录
type AuthQq struct {
	Config *oauth.AuthConfig
	oauth.AuthOauthURL
}

func NewAuthQq(conf *oauth.AuthConfig) *AuthQq {
	auth := oauth.AuthOauthURL{}

	auth.Name = "qq"
	auth.AuthorizeUrl = "https://graph.qq.com/oauth2.0/authorize"
	auth.AccessTokenUrl = "https://graph.qq.com/oauth2.0/token"
	auth.RefreshTokenUrl = "https://graph.qq.com/oauth2.0/token"
	auth.OpenidUrl = "https://graph.qq.com/oauth2.0/me"
	auth.UserInfoUrl = "https://graph.qq.com/user/get_user_info"

	return &AuthQq{
		Config:       conf,
		AuthOauthURL: auth,
	}
}

func (a *AuthQq) GetName() string {
	return a.Name
}

// 获取登录地址
func (a *AuthQq) GetAuthorizeUrl(state string) string {

	url := httpx.NewClient(
		"GET",
		a.AuthorizeUrl,
		httpx.WithParams(map[string]string{
			"client_id":    a.Config.ClientId,
			"redirect_uri": a.Config.RedirectUri,
			//"scope": "contact:user.base:readonly",
			"state":         state,
			"response_type": "code",
		}),
	).EncodeURL()

	return url
}

func (a *AuthQq) GetUserOpenInfo(code string) (resp *oauth.UserResult, err error) {
	token, err := a.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	open, err := a.GetOpenid(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserInfo(token.AccessToken, open.OpenId)
	if err != nil {
		return nil, err
	}

	resp = &oauth.UserResult{
		OpenId:   open.OpenId,
		NickName: user.Nickname,
		Name:     user.Nickname,
		EnName:   user.Nickname,
		Avatar:   user.FigureURLQQ,
		Email:    "",
		Mobile:   "",
	}
	return resp, nil
}

// 获取token
func (a *AuthQq) GetAccessToken(code string) (resp *TokenResult, err error) {

	body, err := httpx.NewClient(
		"GET",
		a.AccessTokenUrl,
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"client_secret": a.Config.ClientSecret,
			"redirect_uri":  a.Config.RedirectUri,
			"code":          code,
			"grant_type":    "authorization_code",
			"fmt":           "json", // 由于历史原因，加上这个参数则返回json格式数据
		}),
	).DoRequest()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 刷新token
func (a *AuthQq) RefreshToken(refreshToken string) (resp *RefreshResult, err error) {

	body, err := httpx.NewClient(
		"GET",
		a.RefreshTokenUrl,
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"client_secret": a.Config.ClientSecret,
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
			"fmt":           "json",
		}),
	).DoRequest()
	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	// 由于QQ的返回值expire_in是字符串，所以不能直接把string解析到int上
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 获取用户openid
func (a *AuthQq) GetOpenid(accessToken string) (resp *OpenResult, err error) {

	body, err := httpx.NewClient(
		"GET",
		a.OpenidUrl,
		httpx.WithParams(map[string]string{
			"access_token": accessToken,
			"fmt":          "json",
		}),
	).DoRequest()
	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// 获取第三方用户信息 https://wiki.connect.qq.com/get_user_info
func (a *AuthQq) GetUserInfo(accessToken string, openId string) (resp *UserResult, err error) {

	body, err := httpx.NewClient(
		"POST",
		a.UserInfoUrl,
		httpx.WithParams(map[string]string{
			"openid":             openId,
			"access_token":       accessToken,
			"oauth_consumer_key": a.Config.ClientId,
		}),
	).DoRequest()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
