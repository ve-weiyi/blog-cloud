package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type UserRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewUserRouter(svcCtx *svctx.ServiceContext) *UserRouter {
	return &UserRouter{
		svcCtx: svcCtx,
	}
}

func (s *UserRouter) Register(r *gin.RouterGroup) {
	// User
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewUserController(s.svcCtx)
		// 获取用户接口权限
		group.GET("/user/get_user_apis", handler.GetUserApis)
		// 获取用户信息
		group.GET("/user/get_user_info", handler.GetUserInfo)
		// 查询用户登录历史
		group.POST("/user/get_user_login_history_list", handler.GetUserLoginHistoryList)
		// 获取用户菜单权限
		group.GET("/user/get_user_menus", handler.GetUserMenus)
		// 获取用户角色
		group.GET("/user/get_user_roles", handler.GetUserRoles)
		// 修改用户信息
		group.POST("/user/update_user_info", handler.UpdateUserInfo)
	}
}
