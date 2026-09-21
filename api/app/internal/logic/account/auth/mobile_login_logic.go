package auth

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/constants/enums"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/authservice"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type MobileLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 手机验证码登录（自动注册）
func NewMobileLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MobileLoginLogic {
	return &MobileLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MobileLoginLogic) MobileLogin(req *types.MobileLoginReq) (resp *types.LoginResp, err error) {
	_, err = l.svcCtx.NotificationService.VerifyMobileCode(l.ctx, &notificationservice.VerifyMobileCodeRequest{
		Mobile: req.Mobile,
		Scene:  enums.CodeSceneMobileLogin,
		Code:   req.Code,
		BizId:  "",
	})
	if err != nil {
		return nil, err
	}

	out, err := l.svcCtx.AuthService.LoginByMobile(l.ctx, &authservice.LoginByMobileRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	return onLogin(l.ctx, l.svcCtx, out)
}
