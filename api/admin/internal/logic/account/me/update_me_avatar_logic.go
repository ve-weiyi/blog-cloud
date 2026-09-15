package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UpdateMeAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户头像
func NewUpdateMeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeAvatarLogic {
	return &UpdateMeAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMeAvatarLogic) UpdateMeAvatar(req *types.UpdateMeAvatarReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UpdateMeAvatar(l.ctx, &userservice.UpdateMeAvatarRequest{
		Avatar: req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
