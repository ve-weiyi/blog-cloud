package svc

import (
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/config"
	"github.com/ve-weiyi/vkit/adapter/oauthx"
)

// NewOAuthProviders 按平台创建 OAuth 提供方。
func NewOAuthProviders(confs map[string]config.OAuthConf) (map[string]oauthx.OAuthProvider, error) {
	providers := make(map[string]oauthx.OAuthProvider, len(confs))
	for platform, c := range confs {
		p, err := oauthx.NewOAuthProvider(&oauthx.OAuthConfig{
			Platform:     platform,
			ClientId:     c.ClientId,
			ClientSecret: c.ClientSecret,
			RedirectUri:  c.RedirectUri,
		})
		if err != nil {
			return nil, err
		}
		providers[platform] = p
	}
	return providers, nil
}
