package auth

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PhoneLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手机登录
func NewPhoneLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneLoginLogic {
	return &PhoneLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneLoginLogic) PhoneLogin(req *types.PhoneLoginReq) (resp *types.LoginResp, err error) {
	in := accountrpc.PhoneLoginReq{
		Phone:      req.Phone,
		VerifyCode: req.VerifyCode,
	}

	out, err := l.svcCtx.AccountRpc.PhoneLogin(l.ctx, &in)
	if err != nil {
		return
	}

	// 登录日志
	_, err = l.svcCtx.SyslogRpc.AddLoginLog(l.ctx, &syslogrpc.AddLoginLogReq{
		UserId:    out.User.UserId,
		LoginType: out.LoginType,
	})

	return onLogin(l.ctx, l.svcCtx, out)
}
