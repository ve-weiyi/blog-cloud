package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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

func (l *GetMeLogic) GetMe(req *types.GetMeReq) (resp *types.UserProfile, err error) {
	out, err := l.svcCtx.UserService.GetMe(l.ctx, &userservice.GetMeRequest{})
	if err != nil {
		return nil, err
	}

	me := out.MeInfo
	return &types.UserProfile{
		UserId:     me.User.UserId,
		Username:   me.User.Username,
		Nickname:   me.User.Nickname,
		Avatar:     me.User.Avatar,
		Email:      me.User.Email,
		Mobile:     me.User.Mobile,
		Status:     me.User.Status,
		CreatedAt:  0,
		UpdatedAt:  0,
		ThirdParty: make([]*types.UserThirdPartyInfo, 0),
		Roles:      make([]string, 0),
		Perms:      make([]string, 0),
	}, nil
}
