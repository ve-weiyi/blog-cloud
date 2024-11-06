package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type MenuController struct {
	svcCtx *svctx.ServiceContext
}

func NewMenuController(svcCtx *svctx.ServiceContext) *MenuController {
	return &MenuController{
		svcCtx: svcCtx,
	}
}

// @Tags		Menu
// @Summary		"创建菜单"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.MenuNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.MenuBackDTO}	"返回信息"
// @Router		/admin_api/v1/menu/add_menu [POST]
func (s *MenuController) AddMenu(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.MenuNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewMenuService(s.svcCtx).AddMenu(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"批量删除菜单"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/menu/batch_delete_menu [DELETE]
func (s *MenuController) BatchDeleteMenu(c *gin.Context) {
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

	data, err := service.NewMenuService(s.svcCtx).BatchDeleteMenu(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"清空菜单列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/menu/clean_menu_list [POST]
func (s *MenuController) CleanMenuList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewMenuService(s.svcCtx).CleanMenuList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"删除菜单"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/menu/delete_menu [DELETE]
func (s *MenuController) DeleteMenu(c *gin.Context) {
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

	data, err := service.NewMenuService(s.svcCtx).DeleteMenu(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"分页获取菜单列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.MenuQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/menu/find_menu_list [POST]
func (s *MenuController) FindMenuList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.MenuQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewMenuService(s.svcCtx).FindMenuList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"同步菜单列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.SyncMenuReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/menu/sync_menu_list [POST]
func (s *MenuController) SyncMenuList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.SyncMenuReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewMenuService(s.svcCtx).SyncMenuList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Menu
// @Summary		"更新菜单"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.MenuNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.MenuBackDTO}	"返回信息"
// @Router		/admin_api/v1/menu/update_menu [PUT]
func (s *MenuController) UpdateMenu(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.MenuNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewMenuService(s.svcCtx).UpdateMenu(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
