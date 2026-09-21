package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type BindMeMobileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定手机号
func NewBindMeMobileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeMobileLogic {
	return &BindMeMobileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindMeMobileLogic) BindMeMobile(req *types.BindMeMobileReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.BindMeMobile(l.ctx, &userservice.BindMeMobileRequest{
		Mobile: req.Mobile,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
