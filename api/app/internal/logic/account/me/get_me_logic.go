package me

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type GetMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取当前用户信息
func NewGetMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeLogic {
	return &GetMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeLogic) GetMe(req *types.GetMeReq) (resp *types.GetMeResp, err error) {
	gr, err := l.svcCtx.UserService.GetMe(l.ctx, &userservice.GetMeRequest{})
	if err != nil {
		return nil, err
	}

	mi := gr.MeInfo
	if mi == nil {
		return nil, fmt.Errorf("meInfo is nil")
	}

	return &types.GetMeResp{
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
