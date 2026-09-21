package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UpdateMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新当前用户信息
func NewUpdateMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeLogic {
	return &UpdateMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMeLogic) UpdateMe(req *types.UpdateMeReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UpdateMe(l.ctx, &userservice.UpdateMeRequest{
		Nickname: req.Nickname,
		Avatar:   &req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
