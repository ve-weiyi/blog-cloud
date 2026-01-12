package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Param		data	body		types.QueryTalkReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/talk/find_talk_list [POST]
func (s *TalkController) FindTalkList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryTalkReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTalkLogic(s.svcCtx).FindTalkList(reqCtx, req)
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
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.Talk}	"返回信息"
// @Router		/blog-api/v1/talk/get_talk [POST]
func (s *TalkController) GetTalk(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTalkLogic(s.svcCtx).GetTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Talk
// @Summary		"点赞说说"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/blog-api/v1/talk/like_talk [PUT]
func (s *TalkController) LikeTalk(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.IdReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewTalkLogic(s.svcCtx).LikeTalk(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
