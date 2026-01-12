package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type RoleLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewRoleLogic(svcCtx *svctx.ServiceContext) *RoleLogic {
	return &RoleLogic{
		svcCtx: svcCtx,
	}
}

// 创建角色
func (s *RoleLogic) AddRole(reqCtx *request.Context, in *types.NewRoleReq) (out *types.RoleBackVO, err error) {
	// todo

	return
}

// 删除角色
func (s *RoleLogic) DeletesRole(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取角色列表
func (s *RoleLogic) FindRoleList(reqCtx *request.Context, in *types.QueryRoleReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 获取角色资源列表
func (s *RoleLogic) FindRoleResources(reqCtx *request.Context, in *types.IdReq) (out *types.RoleResourcesResp, err error) {
	// todo

	return
}

// 更新角色
func (s *RoleLogic) UpdateRole(reqCtx *request.Context, in *types.NewRoleReq) (out *types.RoleBackVO, err error) {
	// todo

	return
}

// 更新角色接口权限
func (s *RoleLogic) UpdateRoleApis(reqCtx *request.Context, in *types.UpdateRoleApisReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 更新角色菜单权限
func (s *RoleLogic) UpdateRoleMenus(reqCtx *request.Context, in *types.UpdateRoleMenusReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
