package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除角色
func (l *DeleteRoleLogic) DeleteRole(in *permissionrpc.DeleteRoleRequest) (*permissionrpc.DeleteRoleResponse, error) {
	if len(in.Ids) == 0 {
		return &permissionrpc.DeleteRoleResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.TUserRoleModel.DeleteBatch(l.ctx, "role_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TRoleModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.DeleteRoleResponse{SuccessCount: rows}, nil
}
