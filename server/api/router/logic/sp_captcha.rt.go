package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type CaptchaRouter struct {
	svcCtx *svc.RouterContext
}

func NewCaptchaRouter(ctx *svc.RouterContext) *CaptchaRouter {
	return &CaptchaRouter{
		svcCtx: ctx,
	}
}

// InitCaptchaRouter 初始化 Captcha 路由信息
func (s *CaptchaRouter) InitCaptchaRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {
	emailPublicRouter := publicRouter.Group("")

	var self = s.svcCtx.CaptchaController
	{
		emailPublicRouter.POST("captcha/email", self.SendCaptchaEmail) // 发送验证码邮件
		emailPublicRouter.POST("captcha/image", self.GetCaptchaImage)  // 获取验证码图片
		emailPublicRouter.POST("captcha/verify", self.VerifyCaptcha)   // 校验验证码
	}
}
