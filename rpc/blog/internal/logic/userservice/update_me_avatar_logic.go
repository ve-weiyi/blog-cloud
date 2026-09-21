package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateMeAvatarLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeAvatarLogic {
	return &UpdateMeAvatarLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新当前用户头像
func (l *UpdateMeAvatarLogic) UpdateMeAvatar(in *userrpc.UpdateMeAvatarRequest) (*userrpc.UpdateMeAvatarResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 更新头像
	_, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, map[string]interface{}{
		"avatar": in.Avatar,
	}, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return &userrpc.UpdateMeAvatarResponse{
		Success: true,
	}, nil
}
