package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type WebsiteController struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsiteController(svcCtx *svctx.ServiceContext) *WebsiteController {
	return &WebsiteController{
		svcCtx: svcCtx,
	}
}

// @Tags		Website
// @Summary		"获取后台首页信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.AdminHomeInfo}	"返回信息"
// @Router		/admin_api/v1/admin [GET]
func (s *WebsiteController) GetAdminHomeInfo(c *gin.Context) {
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

	data, err := service.NewWebsiteService(s.svcCtx).GetAdminHomeInfo(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取关于我的信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.AboutMe}	"返回信息"
// @Router		/admin_api/v1/admin/about_me [GET]
func (s *WebsiteController) GetAboutMe(c *gin.Context) {
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

	data, err := service.NewWebsiteService(s.svcCtx).GetAboutMe(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"更新关于我的信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AboutMe		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/admin/about_me [PUT]
func (s *WebsiteController) UpdateAboutMe(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AboutMe
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewWebsiteService(s.svcCtx).UpdateAboutMe(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取网站配置"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.WebsiteConfig}	"返回信息"
// @Router		/admin_api/v1/admin/get_website_config [GET]
func (s *WebsiteController) GetWebsiteConfig(c *gin.Context) {
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

	data, err := service.NewWebsiteService(s.svcCtx).GetWebsiteConfig(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取服务器信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.Server}	"返回信息"
// @Router		/admin_api/v1/admin/system_state [GET]
func (s *WebsiteController) GetSystemState(c *gin.Context) {
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

	data, err := service.NewWebsiteService(s.svcCtx).GetSystemState(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"更新网站配置"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.WebsiteConfig		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/admin/update_website_config [PUT]
func (s *WebsiteController) UpdateWebsiteConfig(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.WebsiteConfig
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewWebsiteService(s.svcCtx).UpdateWebsiteConfig(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
