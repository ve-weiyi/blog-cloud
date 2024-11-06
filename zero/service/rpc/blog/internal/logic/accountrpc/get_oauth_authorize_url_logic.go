package accountrpclogic

import (
	"context"
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/oauth"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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
func (l *GetOauthAuthorizeUrlLogic) GetOauthAuthorizeUrl(in *accountrpc.OauthLoginReq) (*accountrpc.OauthLoginUrlResp, error) {
	var auth oauth.Oauth
	for platform, v := range l.svcCtx.Oauth {
		if platform == in.Platform {
			auth = v
		}
	}

	if auth == nil {
		return nil, fmt.Errorf("platform %s is not support", in.Platform)
	}

	resp := &accountrpc.OauthLoginUrlResp{}
	resp.Url = auth.GetAuthorizeUrl(in.State)
	return resp, nil
}
