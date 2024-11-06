package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
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
// @Param		data	body		dto.FriendQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/friend/find_friend_list [POST]
func (s *FriendController) FindFriendList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FriendQuery
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

// @Tags		Friend
// @Summary		"创建友链"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.FriendNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.FriendBackDTO}	"返回信息"
// @Router		/admin_api/v1/friend/add_friend [POST]
func (s *FriendController) AddFriend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FriendNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFriendService(s.svcCtx).AddFriend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Friend
// @Summary		"批量删除友链"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/friend/batch_delete_friend [DELETE]
func (s *FriendController) BatchDeleteFriend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.IdsReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFriendService(s.svcCtx).BatchDeleteFriend(reqCtx, req)
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
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/friend/delete_friend [DELETE]
func (s *FriendController) DeleteFriend(c *gin.Context) {
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

	data, err := service.NewFriendService(s.svcCtx).DeleteFriend(reqCtx, req)
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
// @Param		data	body		dto.FriendNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.FriendBackDTO}	"返回信息"
// @Router		/admin_api/v1/friend/update_friend [PUT]
func (s *FriendController) UpdateFriend(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.FriendNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewFriendService(s.svcCtx).UpdateFriend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
