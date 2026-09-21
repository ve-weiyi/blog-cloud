package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
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
	// 校验图形验证码（一次性，无论对错都会失效）
	ok, err := l.svcCtx.CaptchaStore.VerifyCaptcha(l.ctx, req.CaptchaKey, req.CaptchaCode)
	if err != nil {
		l.Errorf("verify captcha error: %v", err)
		return nil, err
	}
	if !ok {
		return nil, bizerr.NewBizError(bizcode.CodeVerifyCodeError, "验证码错误或已失效")
	}

	out, err := l.svcCtx.AuthService.LoginByPassword(l.ctx, &authservice.LoginByPasswordRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return onLogin(l.ctx, l.svcCtx, out)
}
