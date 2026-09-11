package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type DeactivateAccountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 停用账号（进入冷静期）
func NewDeactivateAccountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeactivateAccountLogic {
	return &DeactivateAccountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeactivateAccountLogic) DeactivateAccount(req *types.DeactivateAccountReq) (resp *types.DeactivateAccountResp, err error) {
	// 调用 RPC 服务停用账号
	rpcResp, err := l.svcCtx.UserService.DeactivateAccount(l.ctx, &userservice.DeactivateAccountRequest{
		Password: req.Password,
		Reason:   req.Reason,
	})
	if err != nil {
		return nil, err
	}

	return &types.DeactivateAccountResp{
		CoolingPeriodDays:  rpcResp.CoolingPeriodDays,
		CanReactivateUntil: rpcResp.CanReactivateUntil,
	}, nil
}
