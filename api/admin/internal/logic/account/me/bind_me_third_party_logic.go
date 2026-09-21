package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type BindMeThirdPartyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 绑定第三方平台账号
func NewBindMeThirdPartyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMeThirdPartyLogic {
	return &BindMeThirdPartyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindMeThirdPartyLogic) BindMeThirdParty(req *types.BindMeThirdPartyReq) (resp *types.EmptyResp, err error) {
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
