package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Param		data	body		types.AlbumQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/album/find_album_list [POST]
func (s *AlbumController) FindAlbumList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.AlbumQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAlbumLogic(s.svcCtx).FindAlbumList(reqCtx, req)
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
// @Param		data	body		types.PhotoQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/album/find_photo_list [POST]
func (s *AlbumController) FindPhotoList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.PhotoQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAlbumLogic(s.svcCtx).FindPhotoList(reqCtx, req)
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
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Album}	"返回信息"
// @Router		/blog-api/v1/album/get_album [POST]
func (s *AlbumController) GetAlbum(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewAlbumLogic(s.svcCtx).GetAlbum(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
