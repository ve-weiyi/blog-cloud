package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type BlogApiController struct {
	svcCtx *svctx.ServiceContext
}

func NewBlogApiController(svcCtx *svctx.ServiceContext) *BlogApiController {
	return &BlogApiController{
		svcCtx: svcCtx,
	}
}

// @Tags		BlogApi
// @Summary		"ping"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.PingReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PingResp}	"返回信息"
// @Router		/blog-api/v1/ping [GET]
func (s *BlogApiController) Ping(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.PingReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewBlogApiLogic(s.svcCtx).Ping(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
