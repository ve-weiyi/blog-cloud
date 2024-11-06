package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AccountRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewAccountRouter(svcCtx *svctx.ServiceContext) *AccountRouter {
	return &AccountRouter{
		svcCtx: svcCtx,
	}
}

func (s *AccountRouter) Register(r *gin.RouterGroup) {
	// Account
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewAccountController(s.svcCtx)
		// 获取用户分布地区
		group.POST("/account/find_account_area_analysis", handler.FindAccountAreaAnalysis)
		// 查询用户列表
		group.POST("/account/find_account_list", handler.FindAccountList)
		// 查询用户登录历史
		group.POST("/account/find_account_login_history_list", handler.FindAccountLoginHistoryList)
		// 查询在线用户列表
		group.POST("/account/find_account_online_list", handler.FindAccountOnlineList)
		// 修改用户角色
		group.POST("/account/update_account_roles", handler.UpdateAccountRoles)
		// 修改用户状态
		group.POST("/account/update_account_status", handler.UpdateAccountStatus)
	}
}
