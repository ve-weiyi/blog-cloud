package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户状态
func (l *UpdateUserStatusLogic) UpdateUserStatus(in *userrpc.UpdateUserStatusRequest) (*userrpc.UpdateUserStatusResponse, error) {
	_, err := l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"status": in.Status,
	}, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}

	return &userrpc.UpdateUserStatusResponse{
		Success: true,
	}, nil
}
