package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type VisitorController struct {
	svcCtx *svctx.ServiceContext
}

func NewVisitorController(svcCtx *svctx.ServiceContext) *VisitorController {
	return &VisitorController{
		svcCtx: svcCtx,
	}
}

// @Tags		Visitor
// @Summary		"分页获取游客列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryVisitorReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/visitor/find_visitor_list [POST]
func (s *VisitorController) FindVisitorList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryVisitorReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewVisitorLogic(s.svcCtx).FindVisitorList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
