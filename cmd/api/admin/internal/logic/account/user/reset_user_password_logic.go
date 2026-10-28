package user

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/userservice"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
)

type ResetUserPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重置用户密码
func NewResetUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetUserPasswordLogic {
	return &ResetUserPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetUserPasswordLogic) ResetUserPassword(req *types.ResetUserPasswordReq) (resp *types.ResetUserPasswordResp, err error) {
	_, err = l.svcCtx.UserService.ResetUserPassword(l.ctx, &userservice.ResetUserPasswordRequest{
		UserId:      req.UserId,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	if err := l.svcCtx.TokenManager.RevokeAll(l.ctx, req.UserId); err != nil {
		return nil, err
	}

	return &types.ResetUserPasswordResp{}, nil
}
