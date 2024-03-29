package logic

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/controller/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/response"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/controller"
)

type PhotoAlbumController struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func NewPhotoAlbumController(svcCtx *svc.ControllerContext) *PhotoAlbumController {
	return &PhotoAlbumController{
		svcCtx:         svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		PhotoAlbum
// @Summary		创建相册
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		entity.PhotoAlbum		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.PhotoAlbum}	"返回信息"
// @Router		/photo_album [post]
func (s *PhotoAlbumController) CreatePhotoAlbum(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var photoAlbum entity.PhotoAlbum
	err = s.ShouldBind(c, &photoAlbum)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.PhotoAlbumService.CreatePhotoAlbum(reqCtx, &photoAlbum)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	PhotoAlbum
// @Summary		更新相册
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	data	body 	 	entity.PhotoAlbum		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.PhotoAlbum}	"返回信息"
// @Router 		/photo_album [put]
func (s *PhotoAlbumController) UpdatePhotoAlbum(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var photoAlbum entity.PhotoAlbum
	err = s.ShouldBind(c, &photoAlbum)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.PhotoAlbumService.UpdatePhotoAlbum(reqCtx, &photoAlbum)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		PhotoAlbum
// @Summary		删除相册
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"PhotoAlbum.id"
// @Success		200		{object}	response.Response{data=any}			"返回信息"
// @Router		/photo_album/{id} [delete]
func (s *PhotoAlbumController) DeletePhotoAlbum(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var id int
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.PhotoAlbumService.DeletePhotoAlbum(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	PhotoAlbum
// @Summary		查询相册
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"PhotoAlbum.id"
// @Success		200		{object}	response.Response{data=entity.PhotoAlbum}	"返回信息"
// @Router 		/photo_album/{id} [get]
func (s *PhotoAlbumController) FindPhotoAlbum(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var id int
	id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.PhotoAlbumService.FindPhotoAlbum(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	PhotoAlbum
// @Summary		批量删除相册
// @Accept 	 	application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data 	body		[]int 						true 	"删除id列表"
// @Success		200		{object}	response.Response{data=response.BatchResult}	"返回信息"
// @Router		/photo_album/batch_delete [delete]
func (s *PhotoAlbumController) DeletePhotoAlbumByIds(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var ids []int
	err = s.ShouldBind(c, &ids)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.PhotoAlbumService.DeletePhotoAlbumByIds(reqCtx, ids)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, response.BatchResult{
		SuccessCount: data,
	})
}

// @Tags 	 	PhotoAlbum
// @Summary		分页获取相册列表
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	page 	body		request.PageQuery 			true 	"分页参数"
// @Success		200		{object}	response.Response{data=response.PageResult{list=[]entity.PhotoAlbum}}	"返回信息"
// @Router		/photo_album/list [post]
func (s *PhotoAlbumController) FindPhotoAlbumList(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var page request.PageQuery
	err = s.ShouldBind(c, &page)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	list, total, err := s.svcCtx.PhotoAlbumService.FindPhotoAlbumList(reqCtx, &page)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, response.PageResult{
		List:     list,
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
	})
}
