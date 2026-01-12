package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
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
// @Summary		"创建友链"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewFriendReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.FriendBackVO}	"返回信息"
// @Router		/admin-api/v1/friend/add_friend [POST]
func (s *FriendController) AddFriend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewFriendReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewFriendLogic(s.svcCtx).AddFriend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Friend
// @Summary		"删除友链"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/friend/deletes_friend [DELETE]
func (s *FriendController) DeletesFriend(c *gin.Context) {
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

	data, err := logic.NewFriendLogic(s.svcCtx).DeletesFriend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Friend
// @Summary		"分页获取友链列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryFriendReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/friend/find_friend_list [POST]
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

// @Tags		Friend
// @Summary		"更新友链"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewFriendReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.FriendBackVO}	"返回信息"
// @Router		/admin-api/v1/friend/update_friend [PUT]
func (s *FriendController) UpdateFriend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewFriendReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewFriendLogic(s.svcCtx).UpdateFriend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
