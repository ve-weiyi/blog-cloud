package me

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
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

func (l *GetMeLogic) GetMe(req *types.GetMeReq) (resp *types.GetMeResp, err error) {
	gr, err := l.svcCtx.UserService.GetMe(l.ctx, &userservice.GetMeRequest{})
	if err != nil {
		return nil, err
	}

	mi := gr.MeInfo
	if mi == nil {
		return nil, fmt.Errorf("meInfo is nil")
	}

	thirdParties := make([]*types.UserThirdPartyInfo, 0, len(mi.ThirdParties))
	for _, tp := range mi.ThirdParties {
		thirdParties = append(thirdParties, &types.UserThirdPartyInfo{
			Platform:  tp.Platform,
			OpenId:    tp.OpenId,
			Nickname:  tp.Nickname,
			Avatar:    tp.Avatar,
			CreatedAt: tp.CreatedAt,
		})
	}

	var info types.UserInfoExt
	if mi.User.Info != nil {
		info.Gender = mi.User.Info.Gender
		info.Intro = mi.User.Info.Intro
		info.Website = mi.User.Info.Website
	}

	return &types.GetMeResp{
		UserId:      mi.User.UserId,
		Username:    mi.User.Username,
		Nickname:    mi.User.Nickname,
		Avatar:      mi.User.Avatar,
		Email:       mi.User.Email,
		Mobile:      mi.User.Mobile,
		Status:      mi.User.Status,
		CreatedAt:   mi.User.CreatedAt,
		ThirdParty:  thirdParties,
		UserInfoExt: info,
	}, nil
}
