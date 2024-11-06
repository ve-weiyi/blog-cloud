package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
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
		group := r.Group("/api/v1")

		handler := controller.NewAuthController(s.svcCtx)
		// 登录
		group.POST("/login", handler.Login)
		// 第三方登录授权地址
		group.POST("/oauth_authorize_url", handler.OauthAuthorizeUrl)
		// 第三方登录
		group.POST("/oauth_login", handler.OauthLogin)
		// 注册
		group.POST("/register", handler.Register)
		// 发送注册账号邮件
		group.POST("/send_register_email", handler.SendRegisterEmail)
		// 重置密码
		group.POST("/user/reset_password", handler.ResetPassword)
		// 发送重置密码邮件
		group.POST("/user/send_reset_email", handler.SendResetEmail)
	}
	// Auth
	// [JwtToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewAuthController(s.svcCtx)
		// 绑定邮箱
		group.POST("/bind_user_email", handler.BindUserEmail)
		// 注销
		group.POST("/logoff", handler.Logoff)
		// 登出
		group.POST("/logout", handler.Logout)
		// 发送绑定邮箱验证码
		group.POST("/send_bind_email", handler.SendBindEmail)
	}
}
