package github

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/httpx"
)

// Github授权登录
// https://docs.github.com/zh/apps/overview
type AuthGithub struct {
	Config *oauth.AuthConfig

	Name string // 第三方名称

	AuthorizeUrl    string // 授权登录URL
	AccessTokenUrl  string // 获得应用访问令牌URL
	RefreshTokenUrl string // 刷新令牌URL
	UserInfoUrl     string // 获取用户信息URL
}

func NewAuthGithub(conf *oauth.AuthConfig) *AuthGithub {
	return &AuthGithub{
		Config:         conf,
		Name:           "feishu",
		AuthorizeUrl:   "https://github.com/login/oauth/authorize",
		AccessTokenUrl: "https://github.com/login/oauth/access_token",
		UserInfoUrl:    "https://api.github.com/user",
	}
}

func (a *AuthGithub) GetName() string {
	return a.Name
}

// 1. 获取第三方登录地址（获取授权码code）
func (a *AuthGithub) GetAuthorizeUrl(state string) string {

	url := httpx.NewClient(
		"GET",
		a.AuthorizeUrl,
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"client_secret": a.Config.ClientSecret,
			"redirect_uri":  a.Config.RedirectUri,
			//"scope": "contact:user.base:readonly",
			"state": state,
		}),
	).EncodeURL()
	return url
}

// 获取用户信息
func (a *AuthGithub) GetUserOpenInfo(code string) (resp *oauth.UserResult, err error) {
	token, err := a.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	resp = &oauth.UserResult{
		OpenId:   strconv.Itoa(user.Id),
		NickName: user.Name,
		Name:     user.Login,
		EnName:   user.Login,
		Avatar:   user.AvatarUrl,
		Email:    "",
		Mobile:   "",
	}
	if user.Email != nil {
		resp.Email = *user.Email
	}
	return resp, nil
}

// 获取用户授权凭证
func (a *AuthGithub) GetAccessToken(code string) (resp *Token, err error) {

	body, err := httpx.NewClient(
		"POST",
		a.AccessTokenUrl,
		httpx.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", code),
			"Content-Type":  "application/json; charset=utf-8",
			"Accept":        "application/json",
		}),
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"client_secret": a.Config.ClientSecret,
			"code":          code,
			"redirect_uri":  a.Config.RedirectUri,
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

func (a *AuthGithub) GetUserInfo(accessToken string) (resp *Userinfo, err error) {

	body, err := httpx.NewClient(
		"GET",
		a.UserInfoUrl,
		httpx.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			"Content-Type":  "application/json; charset=utf-8",
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
