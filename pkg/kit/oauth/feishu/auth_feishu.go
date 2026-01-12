package feishu

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/httpx"
)

// Feishu授权登录
type AuthFeishu struct {
	Config *oauth.OauthConfig

	Name string // 第三方名称

	AuthorizeUrl         string // 授权登录URL
	AppAccessTokenUrl    string // 获得应用访问令牌URL
	TenantAccessTokenUrl string // 获取租户授权凭证URL
	UserAccessTokenUrl   string // 获得用户授权凭证URL
	RefreshTokenUrl      string // 刷新令牌URL
	UserInfoUrl          string // 获取用户信息URL
}

func NewAuthFeishu(conf *oauth.OauthConfig) *AuthFeishu {
	return &AuthFeishu{
		Config:               conf,
		Name:                 "feishu",
		AuthorizeUrl:         "https://open.feishu.cn/open-apis/authen/v1/authorize",
		AppAccessTokenUrl:    "https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal",
		TenantAccessTokenUrl: "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal",
		UserAccessTokenUrl:   "https://open.feishu.cn/open-apis/authen/v1/oidc/access_token",
		RefreshTokenUrl:      "https://open.feishu.cn/open-apis/authen/v1/oidc/refresh_access_token",
		UserInfoUrl:          "https://open.feishu.cn/open-apis/authen/v1/user_info",
	}
}

func (a *AuthFeishu) GetName() string {
	return a.Name
}

// 获取登录地址（获取授权码code）
func (a *AuthFeishu) GetAuthLoginUrl(state string) string {

	url := httpx.NewRequest(
		"GET",
		a.AuthorizeUrl,
		httpx.WithParams(map[string]string{
			"app_id":       a.Config.ClientId,
			"redirect_uri": a.Config.RedirectUri,
			"scope":        "contact:user.base:readonly",
			"state":        state,
		}),
	).EncodeURL()
	return url
}

// 获取用户信息
func (a *AuthFeishu) GetAuthUserInfo(code string) (resp *oauth.UserResult, err error) {
	token, err := a.GetUserAccessToken(code)
	if err != nil {
		return nil, err
	}

	out, err := a.GetUserInfo(token.Data.AccessToken)
	if err != nil {
		return nil, err
	}

	user := out.Data

	resp = &oauth.UserResult{
		OpenId:   user.OpenId,
		NickName: user.EnName,
		Name:     user.Name,
		EnName:   user.EnName,
		Avatar:   user.AvatarUrl,
		Email:    user.Email,
		Mobile:   user.Mobile,
	}

	return resp, nil
}

// 获取应用授权凭证
func (a *AuthFeishu) GetAppAccessToken() (resp *AppTokenResp, err error) {

	body, err := httpx.NewRequest(
		"POST",
		a.AppAccessTokenUrl,
		httpx.WithParams(map[string]string{
			"app_id":     a.Config.ClientId,
			"app_secret": a.Config.ClientSecret,
		}),
	).Do()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get app access token failed: %s", resp.Msg)
	}

	return resp, nil
}

// 获取租户授权凭证
func (a *AuthFeishu) GetTenantAccessToken() (resp *TenantTokenResp, err error) {

	body, err := httpx.NewRequest(
		"POST",
		a.TenantAccessTokenUrl,
		httpx.WithParams(map[string]string{
			"app_id":     a.Config.ClientId,
			"app_secret": a.Config.ClientSecret,
		}),
	).Do()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get tenant access token failed: %s", resp.Msg)
	}

	return resp, nil
}

// 获取用户授权凭证
func (a *AuthFeishu) GetUserAccessToken(code string) (resp *UserAccessTokenResp, err error) {
	tt, err := a.GetAppAccessToken()
	if err != nil {
		return nil, err
	}

	body, err := httpx.NewRequest(
		"POST",
		a.UserAccessTokenUrl,
		httpx.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tt.AppAccessToken),
			"Content-Type":  "application/json; charset=utf-8",
		}),
		httpx.WithBodyJson(map[string]any{
			"grant_type": "authorization_code",
			"code":       code,
		}),
	).Do()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get user access token failed: %s", resp.Msg)
	}

	return resp, nil
}

// 获取用户信息
func (a *AuthFeishu) RefreshAccessToken(refreshToken string) (resp *UserAccessTokenResp, err error) {
	tt, err := a.GetAppAccessToken()
	if err != nil {
		return nil, err
	}

	body, err := httpx.NewRequest(
		"POST",
		a.RefreshTokenUrl,
		httpx.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", tt.AppAccessToken),
			"Content-Type":  "application/json; charset=utf-8",
		}),
		httpx.WithBodyJson(map[string]interface{}{
			"grant_type":    "refresh_token",
			"refresh_token": refreshToken,
		}),
	).Do()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("refresh access token failed: %s", resp.Msg)
	}

	return resp, nil
}

func (a *AuthFeishu) GetUserInfo(accessToken string) (resp *UserInfoResp, err error) {

	body, err := httpx.NewRequest(
		"GET",
		a.UserInfoUrl,
		httpx.WithHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", accessToken),
			"Content-Type":  "application/json; charset=utf-8",
		}),
	).Do()

	if err != nil {
		return nil, err
	}

	log.Println("body:", string(body))

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("get user access token failed: %s", resp.Msg)
	}

	return resp, nil
}
