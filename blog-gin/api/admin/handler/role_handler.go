package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type RoleController struct {
	svcCtx *svctx.ServiceContext
}

func NewRoleController(svcCtx *svctx.ServiceContext) *RoleController {
	return &RoleController{
		svcCtx: svcCtx,
	}
}

// @Tags		Role
// @Summary		"创建角色"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.RoleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.RoleBackVO}	"返回信息"
// @Router		/admin-api/v1/role/add_role [POST]
func (s *RoleController) AddRole(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RoleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRoleLogic(s.svcCtx).AddRole(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"删除角色"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/role/deletes_role [POST]
func (s *RoleController) DeletesRole(c *gin.Context) {
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

	data, err := logic.NewRoleLogic(s.svcCtx).DeletesRole(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"分页获取角色列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.RoleQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/role/find_role_list [POST]
func (s *RoleController) FindRoleList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RoleQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRoleLogic(s.svcCtx).FindRoleList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"获取角色资源列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.RoleResourcesResp}	"返回信息"
// @Router		/admin-api/v1/role/find_role_resources [POST]
func (s *RoleController) FindRoleResources(c *gin.Context) {
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

	data, err := logic.NewRoleLogic(s.svcCtx).FindRoleResources(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"更新角色"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.RoleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.RoleBackVO}	"返回信息"
// @Router		/admin-api/v1/role/update_role [PUT]
func (s *RoleController) UpdateRole(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.RoleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRoleLogic(s.svcCtx).UpdateRole(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"更新角色接口权限"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateRoleApisReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/role/update_role_apis [POST]
func (s *RoleController) UpdateRoleApis(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateRoleApisReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRoleLogic(s.svcCtx).UpdateRoleApis(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"更新角色菜单权限"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateRoleMenusReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/role/update_role_menus [POST]
func (s *RoleController) UpdateRoleMenus(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateRoleMenusReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewRoleLogic(s.svcCtx).UpdateRoleMenus(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
