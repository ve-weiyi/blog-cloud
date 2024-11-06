package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type BannerController struct {
	svcCtx *svctx.ServiceContext
}

func NewBannerController(svcCtx *svctx.ServiceContext) *BannerController {
	return &BannerController{
		svcCtx: svcCtx,
	}
}

// @Tags		Banner
// @Summary		"创建页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.BannerNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BannerBackDTO}	"返回信息"
// @Router		/admin_api/v1/banner/add_banner [POST]
func (s *BannerController) AddBanner(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.BannerNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewBannerService(s.svcCtx).AddBanner(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Banner
// @Summary		"删除页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/banner/delete_banner [DELETE]
func (s *BannerController) DeleteBanner(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.IdReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewBannerService(s.svcCtx).DeleteBanner(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Banner
// @Summary		"分页获取页面列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.BannerQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/banner/find_banner_list [POST]
func (s *BannerController) FindBannerList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.BannerQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewBannerService(s.svcCtx).FindBannerList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Banner
// @Summary		"更新页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.BannerNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BannerBackDTO}	"返回信息"
// @Router		/admin_api/v1/banner/update_banner [PUT]
func (s *BannerController) UpdateBanner(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.BannerNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewBannerService(s.svcCtx).UpdateBanner(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
