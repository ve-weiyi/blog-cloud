package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type RoleRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewRoleRouter(svcCtx *svctx.ServiceContext) *RoleRouter {
	return &RoleRouter{
		svcCtx: svcCtx,
	}
}

func (s *RoleRouter) Register(r *gin.RouterGroup) {
	// Role
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewRoleController(s.svcCtx)
		// 创建角色
		group.POST("/role/add_role", h.AddRole)
		// 删除角色
		group.DELETE("/role/deletes_role", h.DeletesRole)
		// 分页获取角色列表
		group.POST("/role/find_role_list", h.FindRoleList)
		// 获取角色资源列表
		group.POST("/role/find_role_resources", h.FindRoleResources)
		// 更新角色
		group.PUT("/role/update_role", h.UpdateRole)
		// 更新角色接口权限
		group.PUT("/role/update_role_apis", h.UpdateRoleApis)
		// 更新角色菜单权限
		group.PUT("/role/update_role_menus", h.UpdateRoleMenus)
	}
}
