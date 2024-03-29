package feishu

import (
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/httpx"
)

// Feishu授权登录
type AuthFeishu struct {
	oauth.AuthOauth
}

func NewAuthFeishu(conf *oauth.AuthConfig) *AuthFeishu {
	authRequest := &AuthFeishu{}
	authRequest.Set("", conf)

	authRequest.AuthorizeUrl = "https://passport.feishu.cn/suite/passport/oauth/authorize"
	authRequest.TokenUrl = "https://passport.feishu.cn/suite/passport/oauth/token"
	authRequest.UserInfoUrl = "https://passport.feishu.cn/suite/passport/oauth/userinfo"
	authRequest.RefreshUrl = "https://passport.feishu.cn/suite/passport/oauth/token"

	return authRequest
}

// 获取登录地址
func (a *AuthFeishu) GetRedirectUrl(state string) string {

	url := httpx.NewClient(
		httpx.WithParam("client_id", a.Config.ClientID),
		httpx.WithParam("redirect_uri", a.Config.RedirectUri),
		httpx.WithParam("response_type", "code"),
		httpx.WithParam("state", state),
	).EncodeURL(a.AuthorizeUrl)
	return url
}

// 获取用户信息
func (a *AuthFeishu) GetUserOpenInfo(code string) (resp *oauth.UserResult, err error) {
	token, err := a.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	resp = &oauth.UserResult{
		OpenID:   user.OpenID,
		NickName: user.EnName,
		Name:     user.Name,
		EnName:   user.EnName,
		Avatar:   user.AvatarURL,
		Email:    user.Email,
		Mobile:   user.Mobile,
	}

	return resp, nil
}

// 获取token https://open.weibo.com/apps/2658270041/privilege/oauth
func (a *AuthFeishu) GetAccessToken(code string) (resp *TokenResult, err error) {

	body, err := httpx.NewClient(
		httpx.WithParam("grant_type", "authorization_code"),
		httpx.WithParam("client_id", a.Config.ClientID),
		httpx.WithParam("client_secret", a.Config.ClientSecret),
		httpx.WithParam("redirect_uri", a.Config.RedirectUri),
		httpx.WithParam("code", code),
	).DoRequest("POST", a.TokenUrl)

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

// 获取用户信息
func (a *AuthFeishu) RefreshToken(refreshToken string) (resp *RefreshResult, err error) {

	body, err := httpx.NewClient(
		httpx.WithBodyObject(map[string]interface{}{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}),
	).DoRequest("POST", a.RefreshUrl)

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

func (a *AuthFeishu) GetUserInfo(accessToken string) (resp *UserResult, err error) {

	body, err := httpx.NewClient(
		httpx.WithHeader("Content-Type", "application/json;charset=UTF-8"),
		httpx.WithHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)),
	).DoRequest("GET", a.UserInfoUrl)

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
