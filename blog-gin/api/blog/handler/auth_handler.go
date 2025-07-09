package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Summary		"获取游客身份信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetTouristInfoResp}	"返回信息"
// @Router		/blog-api/v1/get_tourist_info [GET]
func (s *AuthController) GetTouristInfo(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).GetTouristInfo(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"邮箱登录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmailLoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.LoginResp}	"返回信息"
// @Router		/blog-api/v1/email_login [POST]
func (s *AuthController) EmailLogin(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmailLoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).EmailLogin(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"获取验证码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.GetCaptchaCodeReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetCaptchaCodeResp}	"返回信息"
// @Router		/blog-api/v1/get_captcha_code [POST]
func (s *AuthController) GetCaptchaCode(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetCaptchaCodeReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).GetCaptchaCode(reqCtx, req)
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
// @Param		data	body		types.GetOauthAuthorizeUrlReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetOauthAuthorizeUrlResp}	"返回信息"
// @Router		/blog-api/v1/get_oauth_authorize_url [POST]
func (s *AuthController) GetOauthAuthorizeUrl(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetOauthAuthorizeUrlReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).GetOauthAuthorizeUrl(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"登录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.LoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.LoginResp}	"返回信息"
// @Router		/blog-api/v1/login [POST]
func (s *AuthController) Login(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.LoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).Login(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"手机登录"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.PhoneLoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.LoginResp}	"返回信息"
// @Router		/blog-api/v1/phone_login [POST]
func (s *AuthController) PhoneLogin(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.PhoneLoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).PhoneLogin(reqCtx, req)
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
// @Param		data	body		types.RegisterReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/register [POST]
func (s *AuthController) Register(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RegisterReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).Register(reqCtx, req)
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
// @Param		data	body		types.ResetPasswordReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/reset_password [POST]
func (s *AuthController) ResetPassword(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.ResetPasswordReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).ResetPassword(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"发送邮件验证码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.SendEmailVerifyCodeReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/send_email_verify_code [POST]
func (s *AuthController) SendEmailVerifyCode(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.SendEmailVerifyCodeReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).SendEmailVerifyCode(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Auth
// @Summary		"发送手机验证码"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.SendPhoneVerifyCodeReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/send_phone_verify_code [POST]
func (s *AuthController) SendPhoneVerifyCode(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.SendPhoneVerifyCodeReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).SendPhoneVerifyCode(reqCtx, req)
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
// @Param		data	body		types.ThirdLoginReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.LoginResp}	"返回信息"
// @Router		/blog-api/v1/third_login [POST]
func (s *AuthController) ThirdLogin(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.ThirdLoginReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).ThirdLogin(reqCtx, req)
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/logoff [POST]
func (s *AuthController) Logoff(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).Logoff(reqCtx, req)
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/logout [POST]
func (s *AuthController) Logout(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAuthLogic(s.svcCtx).Logout(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
