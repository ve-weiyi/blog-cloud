package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type BindUserPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定手机号
func NewBindUserPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserPhoneLogic {
	return &BindUserPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUserPhoneLogic) BindUserPhone(req *types.BindUserPhoneReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.BindMePhone(l.ctx, &userservice.BindMePhoneRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
