package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type BindUserEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定邮箱
func NewBindUserEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserEmailLogic {
	return &BindUserEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUserEmailLogic) BindUserEmail(req *types.BindUserEmailReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.BindMeEmail(l.ctx, &userservice.BindMeEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
