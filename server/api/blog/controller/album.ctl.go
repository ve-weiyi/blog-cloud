package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AlbumController struct {
	svcCtx *svctx.ServiceContext
}

func NewAlbumController(svcCtx *svctx.ServiceContext) *AlbumController {
	return &AlbumController{
		svcCtx: svcCtx,
	}
}

// @Tags		Album
// @Summary		"获取相册列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.AlbumQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/album/find_album_list [POST]
func (s *AlbumController) FindAlbumList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.AlbumQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAlbumService(s.svcCtx).FindAlbumList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Album
// @Summary		"获取相册下的照片列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.PhotoQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/album/find_photo_list [POST]
func (s *AlbumController) FindPhotoList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.PhotoQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewAlbumService(s.svcCtx).FindPhotoList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Album
// @Summary		"获取相册"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.Album}	"返回信息"
// @Router		/api/v1/album/get_album [POST]
func (s *AlbumController) GetAlbum(c *gin.Context) {
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

	data, err := service.NewAlbumService(s.svcCtx).GetAlbum(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
