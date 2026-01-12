package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
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
// @Summary		"获取博客前台首页信息"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.GetBlogHomeInfoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetBlogHomeInfoResp}	"返回信息"
// @Router		/blog-api/v1/blog [GET]
func (s *WebsiteController) GetBlogHomeInfo(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetBlogHomeInfoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewWebsiteLogic(s.svcCtx).GetBlogHomeInfo(reqCtx, req)
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
// @Param		data	body		types.GetAboutMeReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.GetAboutMeResp}	"返回信息"
// @Router		/blog-api/v1/blog/about_me [GET]
func (s *WebsiteController) GetAboutMe(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.GetAboutMeReq
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
