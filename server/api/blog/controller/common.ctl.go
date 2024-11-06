package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type CommonController struct {
	svcCtx *svctx.ServiceContext
}

func NewCommonController(svcCtx *svctx.ServiceContext) *CommonController {
	return &CommonController{
		svcCtx: svcCtx,
	}
}

// @Tags		Common
// @Summary		"ping"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.PingReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PingResp}	"返回信息"
// @Router		/api/v1/ping [GET]
func (s *CommonController) Ping(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.PingReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewCommonService(s.svcCtx).Ping(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
