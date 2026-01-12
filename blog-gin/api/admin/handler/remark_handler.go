package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Summary		"删除留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/remark/deletes_remark [DELETE]
func (s *RemarkController) DeletesRemark(c *gin.Context) {
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

	data, err := logic.NewRemarkLogic(s.svcCtx).DeletesRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		"分页获取留言列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryRemarkReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/remark/find_remark_list [POST]
func (s *RemarkController) FindRemarkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryRemarkReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRemarkLogic(s.svcCtx).FindRemarkList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Remark
// @Summary		"更新留言状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateRemarkStatusReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/remark/update_remark_status [PUT]
func (s *RemarkController) UpdateRemarkStatus(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateRemarkStatusReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRemarkLogic(s.svcCtx).UpdateRemarkStatus(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
