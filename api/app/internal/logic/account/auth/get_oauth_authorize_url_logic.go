package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userauthservice"
)

type GetOauthAuthorizeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方登录授权地址
func NewGetOauthAuthorizeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthAuthorizeUrlLogic {
	return &GetOauthAuthorizeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOauthAuthorizeUrlLogic) GetOauthAuthorizeUrl(req *types.GetOauthAuthorizeUrlReq) (resp *types.GetOauthAuthorizeUrlResp, err error) {
	out, err := l.svcCtx.UserAuthService.GetOAuthAuthorizeUrl(l.ctx, &userauthservice.GetOAuthAuthorizeUrlRequest{
		Platform: req.Platform,
		State:    req.State,
	})
	if err != nil {
		return nil, err
	}

	return &types.GetOauthAuthorizeUrlResp{
		AuthorizeUrl: out.AuthorizeUrl,
	}, nil
}
