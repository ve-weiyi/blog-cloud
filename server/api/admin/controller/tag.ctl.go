package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
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
// @Param		data	body		dto.TagQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/tag/find_tag_list [POST]
func (s *TagController) FindTagList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TagQuery
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

// @Tags		Tag
// @Summary		"创建标签"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.TagNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.TagBackDTO}	"返回信息"
// @Router		/admin_api/v1/tag/add_tag [POST]
func (s *TagController) AddTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TagNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTagService(s.svcCtx).AddTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Tag
// @Summary		"批量删除标签"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/tag/batch_delete_tag [DELETE]
func (s *TagController) BatchDeleteTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.IdsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTagService(s.svcCtx).BatchDeleteTag(reqCtx, req)
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
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/tag/delete_tag [DELETE]
func (s *TagController) DeleteTag(c *gin.Context) {
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

	data, err := service.NewTagService(s.svcCtx).DeleteTag(reqCtx, req)
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
// @Param		data	body		dto.TagNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.TagBackDTO}	"返回信息"
// @Router		/admin_api/v1/tag/update_tag [PUT]
func (s *TagController) UpdateTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TagNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTagService(s.svcCtx).UpdateTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
