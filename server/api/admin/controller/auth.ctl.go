package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
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
// @Router		/admin_api/v1/login [POST]
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
// @Summary		"登出"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/logout [POST]
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
