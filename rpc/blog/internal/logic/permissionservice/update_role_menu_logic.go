package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type UpdateRoleMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleMenuLogic {
	return &UpdateRoleMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色菜单资源
func (l *UpdateRoleMenuLogic) UpdateRoleMenu(in *permissionrpc.UpdateRoleMenuRequest) (*permissionrpc.UpdateRoleMenuResponse, error) {
	_, err := l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "role_id = ?", in.RoleId)
	if err != nil {
		return nil, err
	}

	if len(in.MenuIds) > 0 {
		var batch []*model.TRoleMenu
		for _, id := range in.MenuIds {
			batch = append(batch, &model.TRoleMenu{
				RoleId: in.RoleId,
				MenuId: id,
			})
		}
		_, err = l.svcCtx.TRoleMenuModel.InsertBatch(l.ctx, batch...)
		if err != nil {
			return nil, err
		}
	}

	return &permissionrpc.UpdateRoleMenuResponse{Success: true}, nil
}
