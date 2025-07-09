package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
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
		group := r.Group("/blog-api/v1")

		h := handler.NewAuthController(s.svcCtx)
		// 获取游客身份信息
		group.GET("/get_tourist_info", h.GetTouristInfo)
	}
	// Auth
	// [TerminalToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)

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
	// [TerminalToken UserToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)
		group.Use(s.svcCtx.UserToken)

		h := handler.NewAuthController(s.svcCtx)
		// 注销
		group.POST("/logoff", h.Logoff)
		// 登出
		group.POST("/logout", h.Logout)
	}
}
