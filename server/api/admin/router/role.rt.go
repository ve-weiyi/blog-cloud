package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewRoleController(s.svcCtx)
		// 创建角色
		group.POST("/role/add_role", handler.AddRole)
		// 批量删除角色
		group.POST("/role/batch_delete_role", handler.BatchDeleteRole)
		// 删除角色
		group.DELETE("/role/delete_role", handler.DeleteRole)
		// 分页获取角色列表
		group.POST("/role/find_role_list", handler.FindRoleList)
		// 获取角色资源列表
		group.POST("/role/find_role_resources", handler.FindRoleResources)
		// 更新角色
		group.PUT("/role/update_role", handler.UpdateRole)
		// 更新角色接口权限
		group.POST("/role/update_role_apis", handler.UpdateRoleApis)
		// 更新角色菜单权限
		group.POST("/role/update_role_menus", handler.UpdateRoleMenus)
	}
}
