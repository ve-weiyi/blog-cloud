package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type UpdateMePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改用户密码
func NewUpdateMePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMePasswordLogic {
	return &UpdateMePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMePasswordLogic) UpdateMePassword(req *types.UpdateMePasswordReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.UserService.UpdateMePassword(l.ctx, &userservice.UpdateMePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
