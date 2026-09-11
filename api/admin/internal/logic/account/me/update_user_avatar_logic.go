package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UpdateUserAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户头像
func NewUpdateUserAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserAvatarLogic {
	return &UpdateUserAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserAvatarLogic) UpdateUserAvatar(req *types.UpdateUserAvatarReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UpdateMeAvatar(l.ctx, &userservice.UpdateMeAvatarRequest{
		Avatar: req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
