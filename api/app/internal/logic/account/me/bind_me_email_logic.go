package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type BindMeEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定邮箱
func NewBindMeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeEmailLogic {
	return &BindMeEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindMeEmailLogic) BindMeEmail(req *types.BindMeEmailReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.BindMeEmail(l.ctx, &userservice.BindMeEmailRequest{
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
