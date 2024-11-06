package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
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
// @Summary		"分页获取页面列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.BannerQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/banner/find_banner_list [POST]
func (s *BannerController) FindBannerList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.BannerQueryReq
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
