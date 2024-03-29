package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/controller/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/controller"
)

type CaptchaController struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func NewCaptchaController(svcCtx *svc.ControllerContext) *CaptchaController {
	return &CaptchaController{
		svcCtx:         svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		Captcha
// @Summary		发送验证码
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.CaptchaEmailReq		true	"请求body"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/captcha/email [post]
func (s *CaptchaController) SendCaptchaEmail(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.CaptchaEmailReq
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.CaptchaService.SendCaptchaEmail(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Captcha
// @Summary		生成验证码
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.CaptchaReq				true	"请求body"
// @Success		200		{object}	response.Response{data=response.CaptchaDTO}	"生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router		/captcha/image [post]
func (s *CaptchaController) GetCaptchaImage(c *gin.Context) {
	err := s.LimitLock(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.CaptchaReq
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.CaptchaService.GetCaptchaImage(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Captcha
// @Summary		检验验证码
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.CaptchaVerifyReq		true	"请求body"
// @Success		200		{object}	response.Response{data=any}	"生成验证码,返回包括随机数id,base64,验证码长度,是否开启验证码"
// @Router		/captcha/verify [post]
func (s *CaptchaController) VerifyCaptcha(c *gin.Context) {
	err := s.LimitLock(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.CaptchaVerifyReq
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.CaptchaService.VerifyImageCaptcha(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}
