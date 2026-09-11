package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userauthservice"
)

type EmailLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 邮箱验证码登录（仅登录）
func NewEmailLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailLoginLogic {
	return &EmailLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmailLoginLogic) EmailLogin(req *types.EmailLoginReq) (resp *types.LoginResp, err error) {
	_, err = l.svcCtx.NotificationService.VerifyEmailCode(l.ctx, &notificationservice.VerifyEmailCodeRequest{
		Email: req.Email,
		Scene: enums.CodeSceneEmailLogin,
		Code:  req.Code,
		BizId: "",
	})
	if err != nil {
		return nil, err
	}

	out, err := l.svcCtx.UserAuthService.LoginByEmail(l.ctx, &userauthservice.LoginByEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	return onLogin(l.ctx, l.svcCtx, out)
}
