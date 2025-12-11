package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/rpcutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/feishu"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/gitee"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/github"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/qq"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth/weibo"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOauthAuthorizeUrlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOauthAuthorizeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthAuthorizeUrlLogic {
	return &GetOauthAuthorizeUrlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取授权地址
func (l *GetOauthAuthorizeUrlLogic) GetOauthAuthorizeUrl(in *accountrpc.GetOauthAuthorizeUrlReq) (*accountrpc.GetOauthAuthorizeUrlResp, error) {
	app, err := rpcutils.GetAppNameFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	auth, err := GetPlatformOauth(l.ctx, l.svcCtx, app, in.Platform)
	if err != nil {
		return nil, err
	}

	resp := &accountrpc.GetOauthAuthorizeUrlResp{}
	resp.AuthorizeUrl = auth.GetAuthLoginUrl(in.State)
	return resp, nil
}

func GetPlatformOauth(ctx context.Context, svcCtx *svc.ServiceContext, app string, platform string) (oauth.Oauth, error) {
	// 获取APP配置
	appPlatformConf, ok := svcCtx.Config.ThirdPartyConf[app]
	if !ok {
		return nil, fmt.Errorf("app-name %s is not support", app)
	}

	// 获取第三方登录配置
	v, ok := appPlatformConf[platform]
	if !ok {
		return nil, fmt.Errorf("platform %s is not support", platform)
	}

	conf := &oauth.OauthConfig{
		ClientId:     v.ClientId,
		ClientSecret: v.ClientSecret,
		RedirectUri:  v.RedirectUri,
	}

	var auth oauth.Oauth

	switch platform {
	case "qq":
		auth = qq.NewAuthQq(conf)
	case "github":
		auth = github.NewAuthGithub(conf)
	case "gitee":
		auth = gitee.NewAuthGitee(conf)
	case "weibo":
		auth = weibo.NewAuthWb(conf)
	case "feishu":
		auth = feishu.NewAuthFeishu(conf)
	default:
		return nil, fmt.Errorf("platform %s is not support", platform)
	}

	return auth, nil

}
