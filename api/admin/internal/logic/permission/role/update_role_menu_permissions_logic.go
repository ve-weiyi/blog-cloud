package role

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/permissionservice"
)

type UpdateRoleMenuPermissionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新角色菜单权限
func NewUpdateRoleMenuPermissionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleMenuPermissionsLogic {
	return &UpdateRoleMenuPermissionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleMenuPermissionsLogic) UpdateRoleMenuPermissions(req *types.UpdateRoleMenuPermissionsReq) (resp *types.EmptyResp, err error) {
	_, err = l.svcCtx.PermissionService.UpdateRoleMenu(l.ctx, &permissionservice.UpdateRoleMenuRequest{
		RoleId:  req.RoleId,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return nil, err
	}

	return &types.EmptyResp{}, nil
}
