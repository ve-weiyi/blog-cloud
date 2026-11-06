package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/userservice"
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

	thirdParties := make([]*types.UserThirdPartyInfo, 0, len(me.ThirdParties))
	for _, tp := range me.ThirdParties {
		thirdParties = append(thirdParties, &types.UserThirdPartyInfo{
			Platform:  tp.Platform,
			OpenId:    tp.OpenId,
			Nickname:  tp.Nickname,
			Avatar:    tp.Avatar,
			CreatedAt: tp.CreatedAt,
		})
	}

	var info types.UserInfoExt
	if me.User.Info != nil {
		info.Gender = me.User.Info.Gender
		info.Intro = me.User.Info.Intro
		info.Website = me.User.Info.Website
	}

	return &types.UserProfile{
		UserId:      me.User.UserId,
		Username:    me.User.Username,
		Nickname:    me.User.Nickname,
		Avatar:      me.User.Avatar,
		Email:       me.User.Email,
		Mobile:      me.User.Mobile,
		Status:      me.User.Status,
		CreatedAt:   me.User.CreatedAt,
		UpdatedAt:   me.User.UpdatedAt,
		ThirdParty:  thirdParties,
		Roles:       make([]string, 0),
		Perms:       make([]string, 0),
		UserInfoExt: info,
	}, nil
}
