package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
)

type BindUserThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定第三方平台账号
func NewBindUserThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindUserThirdPartyLogic {
	return &BindUserThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindUserThirdPartyLogic) BindUserThirdParty(req *types.BindUserThirdPartyReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.BindMeThirdParty(l.ctx, &userservice.BindMeThirdPartyRequest{
		Platform: req.Platform,
		Code:     req.Code,
		State:    &req.State,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
