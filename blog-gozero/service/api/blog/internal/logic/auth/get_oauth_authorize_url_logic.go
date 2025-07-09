package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
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
	in := &accountrpc.GetOauthAuthorizeUrlReq{
		Platform: req.Platform,
		State:    req.State,
	}

	out, err := l.svcCtx.AccountRpc.GetOauthAuthorizeUrl(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.GetOauthAuthorizeUrlResp{
		AuthorizeUrl: out.AuthorizeUrl,
	}
	return resp, nil
}
