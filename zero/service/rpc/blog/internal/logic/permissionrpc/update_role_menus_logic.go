package permissionrpclogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"
)

type UpdateRoleMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleMenusLogic {
	return &UpdateRoleMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新角色菜单
func (l *UpdateRoleMenusLogic) UpdateRoleMenus(in *permissionrpc.UpdateRoleMenusReq) (*permissionrpc.EmptyResp, error) {
	// 删除角色菜单
	_, err := l.svcCtx.TRoleMenuModel.Deletes(l.ctx, "role_id in (?)", in.RoleId)
	if err != nil {
		return nil, err
	}

	var roleMenus []*model.TRoleMenu
	for _, menuId := range in.MenuIds {
		m := &model.TRoleMenu{
			RoleId: in.RoleId,
			MenuId: menuId,
		}

		roleMenus = append(roleMenus, m)
	}

	// 添加角色菜单
	_, err = l.svcCtx.TRoleMenuModel.Inserts(l.ctx, roleMenus...)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.EmptyResp{}, nil
}
