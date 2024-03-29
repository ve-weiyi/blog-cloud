package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type RoleRouter struct {
	svcCtx *svc.RouterContext
}

func NewRoleRouter(svcCtx *svc.RouterContext) *RoleRouter {
	return &RoleRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Role 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *RoleRouter) InitRoleRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.RoleController
	{
		loginRouter.POST("role", handler.CreateRole)       // 新建Role
		loginRouter.PUT("role", handler.UpdateRole)        // 更新Role
		loginRouter.DELETE("role/:id", handler.DeleteRole) // 删除Role
		loginRouter.GET("role/:id", handler.FindRole)      // 查询Role

		loginRouter.DELETE("role/batch_delete", handler.DeleteRoleByIds) // 批量删除Role列表
		loginRouter.POST("role/list", handler.FindRoleList)              // 分页查询Role列表

		loginRouter.POST("role/details_list", handler.FindRoleDetailsList)
		loginRouter.POST("role/update_menus", handler.UpdateRoleMenus)
		loginRouter.POST("role/update_resources", handler.UpdateRoleResources)
	}
}
