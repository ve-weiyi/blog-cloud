package me

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type GetUserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserProfileLogic) GetUserProfile(req *types.GetUserProfileReq) (resp *types.GetUserProfileResp, err error) {
	gr, err := l.svcCtx.UserService.GetMeProfile(l.ctx, &userservice.GetMeProfileRequest{})
	if err != nil {
		return nil, err
	}

	mi := gr.MeInfo
	if mi == nil {
		return nil, fmt.Errorf("meInfo is nil")
	}

	return &types.GetUserProfileResp{
		UserId:    mi.User.UserId,
		Username:  mi.User.Username,
		Nickname:  mi.User.Nickname,
		Avatar:    mi.User.Avatar,
		Email:     mi.User.Email,
		Mobile:    mi.User.Mobile,
		Status:    mi.User.Status,
		CreatedAt: mi.User.CreatedAt,
	}, nil
}
