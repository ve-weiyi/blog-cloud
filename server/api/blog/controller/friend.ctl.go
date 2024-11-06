package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
// @Param		data	body		dto.FriendQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/friend_link/find_friend_list [POST]
func (s *FriendController) FindFriendList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FriendQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFriendService(s.svcCtx).FindFriendList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
