package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewAccountController(s.svcCtx)
		// 查询用户列表
		group.POST("/account/find_account_list", h.FindAccountList)
		// 查询在线用户列表
		group.POST("/account/find_account_online_list", h.FindAccountOnlineList)
		// 修改用户密码
		group.PUT("/account/update_account_password", h.UpdateAccountPassword)
		// 修改用户角色
		group.PUT("/account/update_account_roles", h.UpdateAccountRoles)
		// 修改用户状态
		group.PUT("/account/update_account_status", h.UpdateAccountStatus)
	}
}
