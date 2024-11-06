package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AuthRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewAuthRouter(svcCtx *svctx.ServiceContext) *AuthRouter {
	return &AuthRouter{
		svcCtx: svcCtx,
	}
}

func (s *AuthRouter) Register(r *gin.RouterGroup) {
	// Auth
	// []
	{
		group := r.Group("/admin_api/v1")

		handler := controller.NewAuthController(s.svcCtx)
		// 登录
		group.POST("/login", handler.Login)
	}
	// Auth
	// [JwtToken]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewAuthController(s.svcCtx)
		// 登出
		group.POST("/logout", handler.Logout)
	}
}
