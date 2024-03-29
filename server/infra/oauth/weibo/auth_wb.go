package weibo

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/httpx"
)

// 微博授权登录
type AuthWb struct {
	oauth.AuthOauth
}

func NewAuthWb(conf *oauth.AuthConfig) *AuthWb {
	authRequest := &AuthWb{}
	authRequest.Set("weibo", conf)

	authRequest.AuthorizeUrl = "https://api.weibo.com/oauth2/authorize"
	authRequest.TokenUrl = "https://api.weibo.com/oauth2/access_token"
	authRequest.UserInfoUrl = "https://api.weibo.com/2/users/show.json"

	return authRequest
}

// 获取登录地址
func (a *AuthWb) GetRedirectUrl(state string) string {

	url := httpx.NewClient(
		httpx.WithParam("response_type", "code"),
		httpx.WithParam("client_id", a.Config.ClientID),
		httpx.WithParam("redirect_uri", a.Config.RedirectUri),
		httpx.WithParam("state", state),
	).EncodeURL(a.AuthorizeUrl)

	return url
}

// 获取用户信息
func (a *AuthWb) GetUserOpenInfo(code string) (resp *oauth.UserResult, err error) {
	tk, err := a.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserInfo(tk.AccessToken, tk.Uid)
	if err != nil {
		return nil, err
	}

	resp = &oauth.UserResult{
		OpenID:   strconv.FormatInt(user.ID, 10),
		NickName: user.ScreenName,
		Name:     user.Name,
		Avatar:   user.AvatarLarge,
	}

	return resp, nil
}

// 获取token
func (a *AuthWb) GetAccessToken(code string) (resp *TokenResult, err error) {

	body, err := httpx.NewClient(
		httpx.WithParam("grant_type", "authorization_code"),
		httpx.WithParam("code", code),
		httpx.WithParam("client_id", a.Config.ClientID),
		httpx.WithParam("client_secret", a.Config.ClientSecret),
		httpx.WithParam("redirect_uri", a.Config.RedirectUri),
	).DoRequest("POST", a.TokenUrl)

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))
	var token TokenResult
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// 获取第三方用户信息
func (a *AuthWb) GetUserInfo(accessToken string, openID string) (resp *UserResult, err error) {

	body, err := httpx.NewClient(
		httpx.WithParam("uid", openID),
		httpx.WithParam("access_token", accessToken),
	).DoRequest("GET", a.UserInfoUrl)
	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	user := UserResult{}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
