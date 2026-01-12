package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ThirdLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方登录
func NewThirdLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdLoginLogic {
	return &ThirdLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdLoginLogic) ThirdLogin(req *types.ThirdLoginReq) (resp *types.LoginResp, err error) {
	in := &accountrpc.ThirdLoginReq{
		Platform: req.Platform,
		Code:     req.Code,
	}

	out, err := l.svcCtx.AccountRpc.ThirdLogin(l.ctx, in)
	if err != nil {
		return
	}

	tk, err := createToken(l.ctx, l.svcCtx, out)
	if err != nil {
		return
	}

	// 登录日志
	_, err = l.svcCtx.SyslogRpc.AddLoginLog(l.ctx, &syslogrpc.NewLoginLogReq{
		UserId:    out.UserId,
		LoginType: out.LoginType,
	})

	resp = &types.LoginResp{
		Token: tk,
	}
	return
}
