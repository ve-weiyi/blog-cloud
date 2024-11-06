package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AuthController struct {
	svcCtx *svctx.ServiceContext
}

func NewAuthController(svcCtx *svctx.ServiceContext) *AuthController {
	return &AuthController{
		svcCtx: svcCtx,
	}
}

// @Tags		Auth
// @Summary		"登录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.LoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.LoginResp}	"返回信息"
// @Router		/api/v1/login [POST]
func (s *AuthController) Login(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.LoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).Login(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"第三方登录授权地址"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.OauthLoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.OauthLoginUrlResp}	"返回信息"
// @Router		/api/v1/oauth_authorize_url [POST]
func (s *AuthController) OauthAuthorizeUrl(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.OauthLoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).OauthAuthorizeUrl(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"第三方登录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.OauthLoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.LoginResp}	"返回信息"
// @Router		/api/v1/oauth_login [POST]
func (s *AuthController) OauthLogin(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.OauthLoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).OauthLogin(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"注册"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.RegisterReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/register [POST]
func (s *AuthController) Register(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RegisterReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).Register(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"发送注册账号邮件"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UserEmailReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/send_register_email [POST]
func (s *AuthController) SendRegisterEmail(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UserEmailReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).SendRegisterEmail(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"重置密码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ResetPasswordReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/user/reset_password [POST]
func (s *AuthController) ResetPassword(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ResetPasswordReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).ResetPassword(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"发送重置密码邮件"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UserEmailReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/user/send_reset_email [POST]
func (s *AuthController) SendResetEmail(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UserEmailReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).SendResetEmail(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"绑定邮箱"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.BindUserEmailReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/bind_user_email [POST]
func (s *AuthController) BindUserEmail(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.BindUserEmailReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).BindUserEmail(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"注销"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/logoff [POST]
func (s *AuthController) Logoff(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).Logoff(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"登出"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/logout [POST]
func (s *AuthController) Logout(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).Logout(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"发送绑定邮箱验证码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UserEmailReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/send_bind_email [POST]
func (s *AuthController) SendBindEmail(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UserEmailReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAuthService(s.svcCtx).SendBindEmail(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
