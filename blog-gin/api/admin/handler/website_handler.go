package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.AdminHomeInfo}	"返回信息"
// @Router		/admin-api/v1/admin [GET]
func (s *WebsiteController) GetAdminHomeInfo(c *gin.Context) {
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

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetAdminHomeInfo(reqCtx, req)
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.AboutMeVO}	"返回信息"
// @Router		/admin-api/v1/admin/get_about_me [GET]
func (s *WebsiteController) GetAboutMe(c *gin.Context) {
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

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetAboutMe(reqCtx, req)
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Server}	"返回信息"
// @Router		/admin-api/v1/admin/get_system_state [GET]
func (s *WebsiteController) GetSystemState(c *gin.Context) {
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

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetSystemState(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取用户分布地区"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.GetUserAreaStatsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetUserAreaStatsResp}	"返回信息"
// @Router		/admin-api/v1/admin/get_user_area_stats [POST]
func (s *WebsiteController) GetUserAreaStats(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetUserAreaStatsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetUserAreaStats(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取访客数据分析"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetVisitStatsResp}	"返回信息"
// @Router		/admin-api/v1/admin/get_visit_stats [GET]
func (s *WebsiteController) GetVisitStats(c *gin.Context) {
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

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetVisitStats(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Website
// @Summary		"获取访客数据趋势"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.GetVisitTrendReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetVisitTrendResp}	"返回信息"
// @Router		/admin-api/v1/admin/get_visit_trend [POST]
func (s *WebsiteController) GetVisitTrend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetVisitTrendReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetVisitTrend(reqCtx, req)
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
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.WebsiteConfigVO}	"返回信息"
// @Router		/admin-api/v1/admin/get_website_config [GET]
func (s *WebsiteController) GetWebsiteConfig(c *gin.Context) {
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

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetWebsiteConfig(reqCtx, req)
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
// @Param		data	body		types.AboutMeVO		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/admin/update_about_me [PUT]
func (s *WebsiteController) UpdateAboutMe(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.AboutMeVO
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewWebsiteLogic(s.svcCtx).UpdateAboutMe(reqCtx, req)
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
// @Param		data	body		types.WebsiteConfigVO		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/admin/update_website_config [PUT]
func (s *WebsiteController) UpdateWebsiteConfig(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.WebsiteConfigVO
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewWebsiteLogic(s.svcCtx).UpdateWebsiteConfig(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
