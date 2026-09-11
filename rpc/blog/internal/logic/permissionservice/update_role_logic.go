package permissionservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色
func (l *UpdateRoleLogic) UpdateRole(in *permissionrpc.UpdateRoleRequest) (*permissionrpc.UpdateRoleResponse, error) {
	role := &model.TRole{
		Id:          in.Id,
		ParentId:    in.ParentId,
		RoleKey:     in.RoleKey,
		RoleLabel:   in.RoleLabel,
		RoleComment: in.RoleComment,
		IsDefault:   in.IsDefault,
		Status:      in.Status,
		UpdatedAt:   time.Now(),
	}

	_, err := l.svcCtx.TRoleModel.Update(l.ctx, role)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.UpdateRoleResponse{Success: true}, nil
}
