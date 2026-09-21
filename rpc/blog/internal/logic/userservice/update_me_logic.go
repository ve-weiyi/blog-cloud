package userservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateMeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMeLogic {
	return &UpdateMeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新当前用户基础信息
func (l *UpdateMeLogic) UpdateMe(in *userrpc.UpdateMeRequest) (*userrpc.UpdateMeResponse, error) {
	// 从上下文获取用户ID
	userID, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 构建更新字段
	updateFields := make(map[string]interface{})
	updateFields["nickname"] = in.Nickname
	if in.Avatar != nil {
		updateFields["avatar"] = *in.Avatar
	}

	// 更新用户信息
	_, err = l.svcCtx.TUserModel.UpdateFields(l.ctx, updateFields, "user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	return &userrpc.UpdateMeResponse{
		Success: true,
	}, nil
}
