package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
		group := r.Group("/admin-api/v1")

		h := handler.NewAuthController(s.svcCtx)
		// 获取客户端信息
		group.GET("/get_client_info", h.GetClientInfo)
	}
	// Auth
	// []
	{
		group := r.Group("/admin-api/v1")

		h := handler.NewAuthController(s.svcCtx)
		// 邮箱登录
		group.POST("/email_login", h.EmailLogin)
		// 获取验证码
		group.POST("/get_captcha_code", h.GetCaptchaCode)
		// 第三方登录授权地址
		group.POST("/get_oauth_authorize_url", h.GetOauthAuthorizeUrl)
		// 登录
		group.POST("/login", h.Login)
		// 手机登录
		group.POST("/phone_login", h.PhoneLogin)
		// 刷新token
		group.POST("/refresh_token", h.RefreshToken)
		// 注册
		group.POST("/register", h.Register)
		// 重置密码
		group.POST("/reset_password", h.ResetPassword)
		// 发送邮件验证码
		group.POST("/send_email_verify_code", h.SendEmailVerifyCode)
		// 发送手机验证码
		group.POST("/send_phone_verify_code", h.SendPhoneVerifyCode)
		// 第三方登录
		group.POST("/third_login", h.ThirdLogin)
	}
	// Auth
	// [AdminToken]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)

		h := handler.NewAuthController(s.svcCtx)
		// 注销
		group.POST("/logoff", h.Logoff)
		// 登出
		group.GET("/logout", h.Logout)
	}
}
