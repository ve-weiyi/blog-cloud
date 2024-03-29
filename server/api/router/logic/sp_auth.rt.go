package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type AuthRouter struct {
	svcCtx *svc.RouterContext
}

func NewAuthRouter(ctx *svc.RouterContext) *AuthRouter {
	return &AuthRouter{
		svcCtx: ctx,
	}
}

// InitAuthRouter 初始化 Auth 路由信息
func (s *AuthRouter) InitAuthRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {
	authPublicRouter := publicRouter.Group("")

	var handler = s.svcCtx.AuthController
	{
		authPublicRouter.POST("login", handler.Login)                  // 登录
		authPublicRouter.GET("logout", handler.Logout)                 // 登出
		authPublicRouter.POST("register", handler.Register)            // 注册
		authPublicRouter.POST("register/email", handler.RegisterEmail) // 注册邮件
		authPublicRouter.POST("oauth/login", handler.OauthLogin)       // 使用第三方登录
		authPublicRouter.POST("oauth/url", handler.GetAuthorizeUrl)    // 获取授权地址
	}
}
