package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"
)

type ReactivateAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 恢复账号（冷静期内）
func NewReactivateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReactivateAccountLogic {
	return &ReactivateAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReactivateAccountLogic) ReactivateAccount(req *types.ReactivateAccountReq) (resp *types.ReactivateAccountResp, err error) {
	// 验证验证码（API 层负责）
	if req.Email != "" {
		_, err = l.svcCtx.NotificationService.VerifyEmailCode(l.ctx, &notificationservice.VerifyEmailCodeRequest{
			Email: req.Email,
			Scene: "reactivate",
			Code:  req.VerificationCode,
			BizId: "", // 空值，使用默认规则 scene:email
		})
		if err != nil {
			return nil, err
		}
	} else if req.Mobile != "" {
		_, err = l.svcCtx.NotificationService.VerifyPhoneCode(l.ctx, &notificationservice.VerifyPhoneCodeRequest{
			Phone: req.Mobile,
			Scene: "reactivate",
			Code:  req.VerificationCode,
			BizId: "", // 空值，使用默认规则 scene:mobile
		})
		if err != nil {
			return nil, err
		}
	}

	// 调用 RPC 服务恢复账号
	rpcResp, err := l.svcCtx.UserService.ReactivateAccount(l.ctx, &userservice.ReactivateAccountRequest{
		Email:  &req.Email,
		Mobile: &req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	return &types.ReactivateAccountResp{
		Success: rpcResp.Success,
	}, nil
}
