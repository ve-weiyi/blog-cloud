package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
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
// @Summary		"获取博客前台首页信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.GetBlogHomeInfoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.GetBlogHomeInfoResp}	"返回信息"
// @Router		/api/v1/blog [GET]
func (s *WebsiteController) GetBlogHomeInfo(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.GetBlogHomeInfoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewWebsiteService(s.svcCtx).GetBlogHomeInfo(reqCtx, req)
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
// @Param		data	body		dto.GetAboutMeReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.GetAboutMeResp}	"返回信息"
// @Router		/api/v1/blog/about_me [GET]
func (s *WebsiteController) GetAboutMe(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.GetAboutMeReq
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
