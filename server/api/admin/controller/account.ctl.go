package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AccountController struct {
	svcCtx *svctx.ServiceContext
}

func NewAccountController(svcCtx *svctx.ServiceContext) *AccountController {
	return &AccountController{
		svcCtx: svcCtx,
	}
}

// @Tags		Account
// @Summary		"获取用户分布地区"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AccountQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/account/find_account_area_analysis [POST]
func (s *AccountController) FindAccountAreaAnalysis(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AccountQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).FindAccountAreaAnalysis(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Account
// @Summary		"查询用户列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AccountQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/account/find_account_list [POST]
func (s *AccountController) FindAccountList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AccountQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).FindAccountList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Account
// @Summary		"查询用户登录历史"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AccountQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/account/find_account_login_history_list [POST]
func (s *AccountController) FindAccountLoginHistoryList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AccountQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).FindAccountLoginHistoryList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Account
// @Summary		"查询在线用户列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AccountQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/account/find_account_online_list [POST]
func (s *AccountController) FindAccountOnlineList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AccountQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).FindAccountOnlineList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Account
// @Summary		"修改用户角色"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UpdateAccountRolesReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/account/update_account_roles [POST]
func (s *AccountController) UpdateAccountRoles(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UpdateAccountRolesReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).UpdateAccountRoles(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Account
// @Summary		"修改用户状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.UpdateAccountStatusReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/account/update_account_status [POST]
func (s *AccountController) UpdateAccountStatus(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UpdateAccountStatusReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAccountService(s.svcCtx).UpdateAccountStatus(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
