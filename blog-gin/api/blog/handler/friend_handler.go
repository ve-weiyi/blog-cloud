package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type FriendController struct {
	svcCtx *svctx.ServiceContext
}

func NewFriendController(svcCtx *svctx.ServiceContext) *FriendController {
	return &FriendController{
		svcCtx: svcCtx,
	}
}

// @Tags		Friend
// @Summary		"分页获取友链列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryFriendReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/blog-api/v1/friend/find_friend_list [POST]
func (s *FriendController) FindFriendList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryFriendReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewFriendLogic(s.svcCtx).FindFriendList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
