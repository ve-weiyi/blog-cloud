package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
)

// 分页获取Role记录
func (s *RoleService) FindRoleDetailsList(reqCtx *request.Context, page *request.PageQuery) (list []*response.RoleDetailsDTO, total int64, err error) {

	// 查找角色列表
	roles, total, err := s.FindRoleList(reqCtx, page)
	if err != nil {
		return nil, 0, err
	}
	var roleIds []int
	var menuMap = make(map[int][]int)
	var apiMap = make(map[int][]int)
	for _, role := range roles {
		roleIds = append(roleIds, role.ID)
	}

	// 查找角色菜单
	menus, err := s.svcCtx.RoleMenuRepository.FindALL(reqCtx, "role_id in (?)", roleIds)
	if err != nil {
		return nil, 0, err
	}

	for _, menu := range menus {
		menuMap[menu.RoleID] = append(menuMap[menu.RoleID], menu.MenuID)
	}

	// 查找角色资源
	apis, err := s.svcCtx.RoleApiRepository.FindALL(reqCtx, "role_id in (?)", roleIds)
	if err != nil {
		return nil, 0, err
	}

	for _, api := range apis {
		apiMap[api.RoleID] = append(apiMap[api.RoleID], api.ApiID)
	}

	// 拼装数据
	for _, role := range roles {
		menuIds := menuMap[role.ID]
		apiIds := apiMap[role.ID]

		r := response.RoleDetailsDTO{
			Role:           *role,
			MenuIdList:     menuIds,
			ResourceIdList: apiIds,
		}
		list = append(list, &r)
	}

	return list, total, nil
}

// 设置角色菜单
func (s *RoleService) UpdateRoleMenus(reqCtx *request.Context, req *request.UpdateRoleMenusReq) (data interface{}, err error) {
	// 重置角色菜单权限
	menu, _, err := s.svcCtx.RoleRepository.UpdateRoleMenus(reqCtx, req.RoleId, req.MenuIds)
	if err != nil {
		return nil, err
	}

	return menu, err
}

// 设置角色菜单
func (s *RoleService) UpdateRoleResources(reqCtx *request.Context, req *request.UpdateRoleApisReq) (data interface{}, err error) {
	// 重置角色接口权限
	role, _, err := s.svcCtx.RoleRepository.UpdateRoleResources(reqCtx, req.RoleId, req.ResourceIds)
	if err != nil {
		return nil, err
	}

	//// 查询资源列表
	//resources, err := s.svcCtx.ApiRepository.FindALL(reqCtx, "api_id in (?)", req.ResourceIds)
	//if err != nil {
	//	return nil, err
	//}
	//// 添加角色规则
	//rbac := s.svcCtx.RBAC
	//rbac.DeleteRolePolicy(role.RoleName, role.RoleDomain)
	//rbac.AddRolePolicy(role.RoleName, role.RoleDomain, resources)

	return role, err
}
