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

type RemarkController struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func NewRemarkController(svcCtx *svc.ControllerContext) *RemarkController {
	return &RemarkController{
		svcCtx:         svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		Remark
// @Summary		创建留言
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		entity.Remark		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.Remark}	"返回信息"
// @Router		/remark [post]
func (s *RemarkController) CreateRemark(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var remark entity.Remark
	err = s.ShouldBind(c, &remark)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.RemarkService.CreateRemark(reqCtx, &remark)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	Remark
// @Summary		更新留言
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	data	body 	 	entity.Remark		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.Remark}	"返回信息"
// @Router 		/remark [put]
func (s *RemarkController) UpdateRemark(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var remark entity.Remark
	err = s.ShouldBind(c, &remark)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.RemarkService.UpdateRemark(reqCtx, &remark)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		删除留言
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"Remark.id"
// @Success		200		{object}	response.Response{data=any}			"返回信息"
// @Router		/remark/{id} [delete]
func (s *RemarkController) DeleteRemark(c *gin.Context) {
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

	data, err := s.svcCtx.RemarkService.DeleteRemark(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	Remark
// @Summary		查询留言
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"Remark.id"
// @Success		200		{object}	response.Response{data=entity.Remark}	"返回信息"
// @Router 		/remark/{id} [get]
func (s *RemarkController) FindRemark(c *gin.Context) {
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

	data, err := s.svcCtx.RemarkService.FindRemark(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	Remark
// @Summary		批量删除留言
// @Accept 	 	application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data 	body		[]int 						true 	"删除id列表"
// @Success		200		{object}	response.Response{data=response.BatchResult}	"返回信息"
// @Router		/remark/batch_delete [delete]
func (s *RemarkController) DeleteRemarkByIds(c *gin.Context) {
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

	data, err := s.svcCtx.RemarkService.DeleteRemarkByIds(reqCtx, ids)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, response.BatchResult{
		SuccessCount: data,
	})
}

// @Tags 	 	Remark
// @Summary		分页获取留言列表
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	page 	body		request.PageQuery 			true 	"分页参数"
// @Success		200		{object}	response.Response{data=response.PageResult{list=[]entity.Remark}}	"返回信息"
// @Router		/remark/list [post]
func (s *RemarkController) FindRemarkList(c *gin.Context) {
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

	list, total, err := s.svcCtx.RemarkService.FindRemarkList(reqCtx, &page)
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
