package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TalkController struct {
	svcCtx *svctx.ServiceContext
}

func NewTalkController(svcCtx *svctx.ServiceContext) *TalkController {
	return &TalkController{
		svcCtx: svcCtx,
	}
}

// @Tags		Talk
// @Summary		"分页获取说说列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.TalkQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/talk/find_talk_list [POST]
func (s *TalkController) FindTalkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TalkQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTalkService(s.svcCtx).FindTalkList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Talk
// @Summary		"创建说说"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.TalkNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.TalkBackDTO}	"返回信息"
// @Router		/admin_api/v1/talk/add_talk [POST]
func (s *TalkController) AddTalk(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TalkNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTalkService(s.svcCtx).AddTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Talk
// @Summary		"删除说说"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/talk/delete_talk [DELETE]
func (s *TalkController) DeleteTalk(c *gin.Context) {
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

	data, err := service.NewTalkService(s.svcCtx).DeleteTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Talk
// @Summary		"查询说说"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.TalkBackDTO}	"返回信息"
// @Router		/admin_api/v1/talk/get_talk [POST]
func (s *TalkController) GetTalk(c *gin.Context) {
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

	data, err := service.NewTalkService(s.svcCtx).GetTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Talk
// @Summary		"更新说说"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.TalkNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.TalkBackDTO}	"返回信息"
// @Router		/admin_api/v1/talk/update_talk [PUT]
func (s *TalkController) UpdateTalk(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.TalkNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewTalkService(s.svcCtx).UpdateTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
