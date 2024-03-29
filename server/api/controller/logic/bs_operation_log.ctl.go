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

type OperationLogController struct {
	controller.BaseController
	svcCtx *svc.ControllerContext
}

func NewOperationLogController(svcCtx *svc.ControllerContext) *OperationLogController {
	return &OperationLogController{
		svcCtx:         svcCtx,
		BaseController: controller.NewBaseController(svcCtx),
	}
}

// @Tags		OperationLog
// @Summary		创建操作记录
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data	body		entity.OperationLog		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.OperationLog}	"返回信息"
// @Router		/operation_log [post]
func (s *OperationLogController) CreateOperationLog(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var operationLog entity.OperationLog
	err = s.ShouldBind(c, &operationLog)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.OperationLogService.CreateOperationLog(reqCtx, &operationLog)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	OperationLog
// @Summary		更新操作记录
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	data	body 	 	entity.OperationLog		true	"请求参数"
// @Success		200		{object}	response.Response{data=entity.OperationLog}	"返回信息"
// @Router 		/operation_log [put]
func (s *OperationLogController) UpdateOperationLog(c *gin.Context) {
	reqCtx, err := s.GetRequestContext(c)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	var operationLog entity.OperationLog
	err = s.ShouldBind(c, &operationLog)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	data, err := s.svcCtx.OperationLogService.UpdateOperationLog(reqCtx, &operationLog)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags		OperationLog
// @Summary		删除操作记录
// @Accept		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"OperationLog.id"
// @Success		200		{object}	response.Response{data=any}			"返回信息"
// @Router		/operation_log/{id} [delete]
func (s *OperationLogController) DeleteOperationLog(c *gin.Context) {
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

	data, err := s.svcCtx.OperationLogService.DeleteOperationLog(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	OperationLog
// @Summary		查询操作记录
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	id		path		int							true	"OperationLog.id"
// @Success		200		{object}	response.Response{data=entity.OperationLog}	"返回信息"
// @Router 		/operation_log/{id} [get]
func (s *OperationLogController) FindOperationLog(c *gin.Context) {
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

	data, err := s.svcCtx.OperationLogService.FindOperationLog(reqCtx, id)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, data)
}

// @Tags 	 	OperationLog
// @Summary		批量删除操作记录
// @Accept 	 	application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param		data 	body		[]int 						true 	"删除id列表"
// @Success		200		{object}	response.Response{data=response.BatchResult}	"返回信息"
// @Router		/operation_log/batch_delete [delete]
func (s *OperationLogController) DeleteOperationLogByIds(c *gin.Context) {
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

	data, err := s.svcCtx.OperationLogService.DeleteOperationLogByIds(reqCtx, ids)
	if err != nil {
		s.ResponseError(c, err)
		return
	}

	s.ResponseOk(c, response.BatchResult{
		SuccessCount: data,
	})
}

// @Tags 	 	OperationLog
// @Summary		分页获取操作记录列表
// @Accept 		application/json
// @Produce		application/json
// @Param		token	header		string						false	"token"
// @Param		uid		header		string						false	"uid"
// @Param 	 	page 	body		request.PageQuery 			true 	"分页参数"
// @Success		200		{object}	response.Response{data=response.PageResult{list=[]entity.OperationLog}}	"返回信息"
// @Router		/operation_log/list [post]
func (s *OperationLogController) FindOperationLogList(c *gin.Context) {
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

	list, total, err := s.svcCtx.OperationLogService.FindOperationLogList(reqCtx, &page)
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
