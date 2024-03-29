package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/controller/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/controller"
)

type AuthController struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func NewAuthController(svcCtx *svc.ControllerContext) *AuthController {
	return &AuthController{
		svcCtx:         svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		Auth
// @Summary		登录
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.UserReq				true	"创建权限认证"
// @Success		200		{object}	response.Response{data=response.Login}	"返回信息"
// @Router		/login [post]
func (s *AuthController) Login(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var user request.UserReq
	err = s.ShouldBind(c, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.Login(reqCtx, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		登出
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/logout [get]
func (s *AuthController) Logout(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.Log.Println("Logout")
	data, err := s.svcCtx.AuthService.Logout(reqCtx, nil)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		注销
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/logoff [post]
func (s *AuthController) Logoff(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.Logoff(reqCtx, nil)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		注册
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.UserReq				true	"请求body"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/register [post]
func (s *AuthController) Register(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var user request.UserReq
	err = s.ShouldBind(c, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.Register(reqCtx, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		发送注册邮件
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.UserEmail			true	"请求body"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/register/email [post]
func (s *AuthController) RegisterEmail(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.UserEmail
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.SendRegisterEmail(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		发送忘记密码邮件
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.UserEmail			true	"请求参数"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/forget/password [post]
func (s *UserController) ForgetPasswordEmail(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var user request.UserEmail
	err = s.ShouldBind(c, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.UserService.SendForgetPwdEmail(reqCtx, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		重置密码
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		request.ResetPasswordReq	true	"请求参数"
// @Success		200		{object}	response.Response{data=any}	"返回信息"
// @Router		/forget/reset_password [post]
func (s *UserController) ResetPassword(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var user request.ResetPasswordReq
	err = s.ShouldBind(c, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.UserService.ResetPassword(reqCtx, &user)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		获取授权地址
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string											false	"token"
// @Param		uid		header		string											false	"uid"
// @Param		data	body		request.OauthLoginReq							true	"请求body"
// @Success		200		{object}	response.Response{data=response.Login}	"返回信息"
// @Router		/oauth/login [post]
func (s *AuthController) OauthLogin(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.OauthLoginReq
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.OauthLogin(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		获取授权地址
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string											false	"token"
// @Param		uid		header		string											false	"uid"
// @Param		data	body		request.OauthLoginReq							true	"请求body"
// @Success		200		{object}	response.Response{data=response.OauthLoginUrl}	"返回信息"
// @Router		/oauth/url [post]
func (s *AuthController) GetAuthorizeUrl(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var req request.OauthLoginReq
	err = s.ShouldBind(c, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.AuthService.GetAuthorizeUrl(reqCtx, &req)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}
