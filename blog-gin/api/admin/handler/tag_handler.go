package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Summary		"创建标签"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.TagNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.TagBackVO}	"返回信息"
// @Router		/admin-api/v1/tag/add_tag [POST]
func (s *TagController) AddTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.TagNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTagLogic(s.svcCtx).AddTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Tag
// @Summary		"删除标签"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/tag/deletes_tag [DELETE]
func (s *TagController) DeletesTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTagLogic(s.svcCtx).DeletesTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Tag
// @Summary		"分页获取标签列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.TagQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/tag/find_tag_list [POST]
func (s *TagController) FindTagList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.TagQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTagLogic(s.svcCtx).FindTagList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Tag
// @Summary		"更新标签"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.TagNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.TagBackVO}	"返回信息"
// @Router		/admin-api/v1/tag/update_tag [PUT]
func (s *TagController) UpdateTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.TagNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTagLogic(s.svcCtx).UpdateTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
