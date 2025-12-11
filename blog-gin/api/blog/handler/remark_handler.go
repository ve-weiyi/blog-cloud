package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
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
// @Summary		"分页获取留言列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.RemarkQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/remark/find_remark_list [POST]
func (s *RemarkController) FindRemarkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RemarkQueryReq
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
// @Summary		"创建留言"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.RemarkNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Remark}	"返回信息"
// @Router		/blog-api/v1/remark/add_remark [POST]
func (s *RemarkController) AddRemark(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RemarkNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRemarkLogic(s.svcCtx).AddRemark(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
