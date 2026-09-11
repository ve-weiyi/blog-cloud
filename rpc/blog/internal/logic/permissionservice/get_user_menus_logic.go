package permissionservicelogic

import (
	"context"
	"slices"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetUserMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserMenusLogic {
	return &GetUserMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询用户菜单权限
func (l *GetUserMenusLogic) GetUserMenus(in *permissionrpc.GetUserMenusRequest) (*permissionrpc.GetUserMenusResponse, error) {
	urs, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id = ?", in.UserId)
	if err != nil {
		return nil, err
	}
	var roleIds []int64
	for _, v := range urs {
		roleIds = append(roleIds, v.RoleId)
	}
	if len(roleIds) == 0 {
		return &permissionrpc.GetUserMenusResponse{}, nil
	}

	rLinks, err := l.svcCtx.TRoleMenuModel.FindALL(l.ctx, "role_id in (?)", roleIds)
	if err != nil {
		return nil, err
	}
	var menuIds []int64
	for _, v := range rLinks {
		if !slices.Contains(menuIds, v.MenuId) {
			menuIds = append(menuIds, v.MenuId)
		}
	}
	if len(menuIds) == 0 {
		return &permissionrpc.GetUserMenusResponse{}, nil
	}

	menus, err := l.svcCtx.TMenuModel.FindALL(l.ctx, "id in (?)", menuIds)
	if err != nil {
		return nil, err
	}

	out := &permissionrpc.GetUserMenusResponse{}
	for _, m := range menus {
		out.List = append(out.List, convertMenuOut(m))
	}

	return out, nil
}
