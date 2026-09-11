package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type UpdateUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新当前用户信息
func NewUpdateUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileLogic {
	return &UpdateUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserProfileLogic) UpdateUserProfile(req *types.UpdateUserProfileReq) (resp *types.UpdateUserProfileResp, err error) {
	_, err = l.svcCtx.UserService.UpdateMe(l.ctx, &userservice.UpdateMeRequest{
		Nickname: req.Nickname,
		Avatar:   &req.Avatar,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateUserProfileResp{}, nil
}
