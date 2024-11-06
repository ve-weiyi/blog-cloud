package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type PhotoController struct {
	svcCtx *svctx.ServiceContext
}

func NewPhotoController(svcCtx *svctx.ServiceContext) *PhotoController {
	return &PhotoController{
		svcCtx: svcCtx,
	}
}

// @Tags		Photo
// @Summary		"分页获取照片列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.PhotoQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/photo/find_photo_list [POST]
func (s *PhotoController) FindPhotoList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.PhotoQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewPhotoService(s.svcCtx).FindPhotoList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"批量删除照片"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/album/batch_delete_photo [DELETE]
func (s *PhotoController) BatchDeletePhoto(c *gin.Context) {
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

	data, err := service.NewPhotoService(s.svcCtx).BatchDeletePhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"创建照片"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.PhotoNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PhotoBackDTO}	"返回信息"
// @Router		/admin_api/v1/photo/add_photo [POST]
func (s *PhotoController) AddPhoto(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.PhotoNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewPhotoService(s.svcCtx).AddPhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"删除照片"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/photo/delete_photo [DELETE]
func (s *PhotoController) DeletePhoto(c *gin.Context) {
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

	data, err := service.NewPhotoService(s.svcCtx).DeletePhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"更新照片"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.PhotoNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PhotoBackDTO}	"返回信息"
// @Router		/admin_api/v1/photo/update_photo [PUT]
func (s *PhotoController) UpdatePhoto(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.PhotoNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewPhotoService(s.svcCtx).UpdatePhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
