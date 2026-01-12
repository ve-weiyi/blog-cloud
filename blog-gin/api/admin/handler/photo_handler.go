package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Summary		"创建照片"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewPhotoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PhotoBackVO}	"返回信息"
// @Router		/admin-api/v1/photo/add_photo [POST]
func (s *PhotoController) AddPhoto(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewPhotoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPhotoLogic(s.svcCtx).AddPhoto(reqCtx, req)
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
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/photo/deletes_photo [DELETE]
func (s *PhotoController) DeletesPhoto(c *gin.Context) {
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

	data, err := logic.NewPhotoLogic(s.svcCtx).DeletesPhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"分页获取照片列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryPhotoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/photo/find_photo_list [POST]
func (s *PhotoController) FindPhotoList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryPhotoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPhotoLogic(s.svcCtx).FindPhotoList(reqCtx, req)
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
// @Param		data	body		types.NewPhotoReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PhotoBackVO}	"返回信息"
// @Router		/admin-api/v1/photo/update_photo [PUT]
func (s *PhotoController) UpdatePhoto(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewPhotoReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPhotoLogic(s.svcCtx).UpdatePhoto(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Photo
// @Summary		"更新照片删除状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdatePhotoDeleteReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/photo/update_photo_delete [PUT]
func (s *PhotoController) UpdatePhotoDelete(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdatePhotoDeleteReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPhotoLogic(s.svcCtx).UpdatePhotoDelete(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
