package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OauthAuthorizeUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方登录授权地址
func NewOauthAuthorizeUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthAuthorizeUrlLogic {
	return &OauthAuthorizeUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OauthAuthorizeUrlLogic) OauthAuthorizeUrl(req *types.OauthLoginReq) (resp *types.OauthLoginUrlResp, err error) {
	in := &accountrpc.OauthLoginReq{
		Platform: req.Platform,
		Code:     req.Code,
		State:    req.State,
	}

	out, err := l.svcCtx.AccountRpc.GetOauthAuthorizeUrl(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.OauthLoginUrlResp{
		Url: out.Url,
	}
	return resp, nil
}
