package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UnbindUserThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 解绑第三方平台账号
func NewUnbindUserThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindUserThirdPartyLogic {
	return &UnbindUserThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnbindUserThirdPartyLogic) UnbindUserThirdParty(req *types.UnbindUserThirdPartyReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UnbindMeThirdParty(l.ctx, &userservice.UnbindMeThirdPartyRequest{
		Platform: req.Platform,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
