package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
)

type PasswordLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 密码登录（账号/手机号/邮箱）
func NewPasswordLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PasswordLoginLogic {
	return &PasswordLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PasswordLoginLogic) PasswordLogin(req *types.PasswordLoginReq) (resp *types.LoginResp, err error) {
	out, err := l.svcCtx.AuthService.LoginByPassword(l.ctx, &authservice.LoginByPasswordRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return onLogin(l.ctx, l.svcCtx, out)
}
