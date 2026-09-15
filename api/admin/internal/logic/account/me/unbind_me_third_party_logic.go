package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UnbindMeThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 解绑第三方平台账号
func NewUnbindMeThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMeThirdPartyLogic {
	return &UnbindMeThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindMeThirdPartyLogic) UnbindMeThirdParty(req *types.UnbindMeThirdPartyReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UnbindMeThirdParty(l.ctx, &userservice.UnbindMeThirdPartyRequest{
		Platform: req.Platform,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
