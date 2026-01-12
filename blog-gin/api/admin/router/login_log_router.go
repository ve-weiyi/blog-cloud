package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type LoginLogRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewLoginLogRouter(svcCtx *svctx.ServiceContext) *LoginLogRouter {
	return &LoginLogRouter{
		svcCtx: svcCtx,
	}
}

func (s *LoginLogRouter) Register(r *gin.RouterGroup) {
	// LoginLog
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewLoginLogController(s.svcCtx)
		// 查询登录日志
		group.POST("/file_log/find_login_log_list", h.FindLoginLogList)
		// 删除登录日志
		group.DELETE("/login_log/deletes_login_log", h.DeletesLoginLog)
	}
}
