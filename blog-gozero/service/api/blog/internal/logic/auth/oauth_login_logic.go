package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OauthLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方登录
func NewOauthLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OauthLoginLogic {
	return &OauthLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OauthLoginLogic) OauthLogin(req *types.OauthLoginReq) (resp *types.LoginResp, err error) {
	in := accountrpc.OauthLoginReq{
		Platform: req.Platform,
		Code:     req.Code,
		State:    req.State,
	}

	out, err := l.svcCtx.AccountRpc.OauthLogin(l.ctx, &in)
	if err != nil {
		return
	}

	tk, err := createToken(l.ctx, l.svcCtx, out)
	if err != nil {
		return
	}

	resp = &types.LoginResp{
		Token: tk,
	}
	return
}
