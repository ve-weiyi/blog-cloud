package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TagController struct {
	svcCtx *svctx.ServiceContext
}

func NewTagController(svcCtx *svctx.ServiceContext) *TagController {
	return &TagController{
		svcCtx: svcCtx,
	}
}

// @Tags		Tag
// @Summary		"分页获取标签列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.TagQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/tag/find_tag_list [POST]
func (s *TagController) FindTagList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TagQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTagService(s.svcCtx).FindTagList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
