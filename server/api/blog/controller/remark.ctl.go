package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
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
// @Param		data	body		dto.RemarkQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/remark/find_remark_list [POST]
func (s *RemarkController) FindRemarkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RemarkQueryReq
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
// @Summary		"创建留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.RemarkNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.Remark}	"返回信息"
// @Router		/api/v1/remark/add_remark [POST]
func (s *RemarkController) AddRemark(c *gin.Context) {
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

	data, err := service.NewRemarkService(s.svcCtx).AddRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
