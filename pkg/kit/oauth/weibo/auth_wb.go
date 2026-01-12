package weibo

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oauth"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/httpx"
)

// 微博授权登录
type AuthWb struct {
	Config *oauth.OauthConfig

	Name           string // 第三方名称
	AuthorizeUrl   string // 授权登录URL
	AccessTokenUrl string // 获得访问令牌URL
	UserInfoUrl    string // 获取用户信息URL
}

func NewAuthWb(conf *oauth.OauthConfig) *AuthWb {
	return &AuthWb{
		Config:         conf,
		Name:           "weibo",
		AuthorizeUrl:   "https://api.weibo.com/oauth/authorize",
		AccessTokenUrl: "https://api.weibo.com/oauth/access_token",
		UserInfoUrl:    "https://api.weibo.com/2/users/show.json",
	}
}

func (a *AuthWb) GetName() string {
	return a.Name
}

// 获取登录地址
func (a *AuthWb) GetAuthLoginUrl(state string) string {

	url := httpx.NewRequest(
		"GET",
		a.AuthorizeUrl,
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"redirect_uri":  a.Config.RedirectUri,
			"state":         state,
			"response_type": "code",
		}),
	).EncodeURL()

	return url
}

// 获取用户信息
func (a *AuthWb) GetAuthUserInfo(code string) (resp *oauth.UserResult, err error) {
	tk, err := a.GetAccessToken(code)
	if err != nil {
		return nil, err
	}

	user, err := a.GetUserInfo(tk.AccessToken, tk.Uid)
	if err != nil {
		return nil, err
	}

	resp = &oauth.UserResult{
		OpenId:   strconv.FormatInt(user.Id, 10),
		NickName: user.ScreenName,
		Name:     user.Name,
		Avatar:   user.AvatarLarge,
	}

	return resp, nil
}

// 获取token
func (a *AuthWb) GetAccessToken(code string) (resp *TokenResult, err error) {

	body, err := httpx.NewRequest(
		"POST",
		a.AccessTokenUrl,
		httpx.WithParams(map[string]string{
			"client_id":     a.Config.ClientId,
			"client_secret": a.Config.ClientSecret,
			"redirect_uri":  a.Config.RedirectUri,
			"code":          code,
			"grant_type":    "authorization_code",
		}),
	).Do()

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
func (a *AuthWb) GetUserInfo(accessToken string, openId string) (resp *UserResult, err error) {

	body, err := httpx.NewRequest(
		"GET",
		a.UserInfoUrl,
		httpx.WithParams(map[string]string{
			"uid":          openId,
			"access_token": accessToken,
		}),
	).Do()
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
