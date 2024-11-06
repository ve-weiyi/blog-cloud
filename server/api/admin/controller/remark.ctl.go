package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type RemarkController struct {
	svcCtx *svctx.ServiceContext
}

func NewRemarkController(svcCtx *svctx.ServiceContext) *RemarkController {
	return &RemarkController{
		svcCtx: svcCtx,
	}
}

// @Tags		Remark
// @Summary		"分页获取留言列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.RemarkQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/remark/find_remark_list [POST]
func (s *RemarkController) FindRemarkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RemarkQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRemarkService(s.svcCtx).FindRemarkList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		"批量删除留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/remark/batch_delete_remark [DELETE]
func (s *RemarkController) BatchDeleteRemark(c *gin.Context) {
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

	data, err := service.NewRemarkService(s.svcCtx).BatchDeleteRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		"删除留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/remark/delete_remark [DELETE]
func (s *RemarkController) DeleteRemark(c *gin.Context) {
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

	data, err := service.NewRemarkService(s.svcCtx).DeleteRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		"更新留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.RemarkNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.RemarkBackDTO}	"返回信息"
// @Router		/admin_api/v1/remark/update_remark [PUT]
func (s *RemarkController) UpdateRemark(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RemarkNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRemarkService(s.svcCtx).UpdateRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
