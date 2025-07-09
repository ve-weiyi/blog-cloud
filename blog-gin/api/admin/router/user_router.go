package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewUserController(s.svcCtx)
		// 删除用户绑定第三方平台账号
		group.POST("/user/delete_user_bind_third_party", h.DeleteUserBindThirdParty)
		// 获取用户接口权限
		group.GET("/user/get_user_apis", h.GetUserApis)
		// 获取用户信息
		group.GET("/user/get_user_info", h.GetUserInfo)
		// 查询用户登录历史
		group.POST("/user/get_user_login_history_list", h.GetUserLoginHistoryList)
		// 获取用户菜单权限
		group.GET("/user/get_user_menus", h.GetUserMenus)
		// 获取用户角色
		group.GET("/user/get_user_roles", h.GetUserRoles)
		// 修改用户头像
		group.POST("/user/update_user_avatar", h.UpdateUserAvatar)
		// 修改用户绑定邮箱
		group.POST("/user/update_user_bind_email", h.UpdateUserBindEmail)
		// 修改用户绑定手机号
		group.POST("/user/update_user_bind_phone", h.UpdateUserBindPhone)
		// 修改用户绑定第三方平台账号
		group.POST("/user/update_user_bind_third_party", h.UpdateUserBindThirdParty)
		// 修改用户信息
		group.POST("/user/update_user_info", h.UpdateUserInfo)
		// 修改用户密码
		group.POST("/user/update_user_password", h.UpdateUserPassword)
	}
}
