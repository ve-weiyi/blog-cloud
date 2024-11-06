package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
// @Param		data	body		dto.RoleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.RoleBackDTO}	"返回信息"
// @Router		/admin_api/v1/role/add_role [POST]
func (s *RoleController) AddRole(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RoleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRoleService(s.svcCtx).AddRole(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Role
// @Summary		"批量删除角色"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/role/batch_delete_role [POST]
func (s *RoleController) BatchDeleteRole(c *gin.Context) {
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

	data, err := service.NewRoleService(s.svcCtx).BatchDeleteRole(reqCtx, req)
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
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/role/delete_role [DELETE]
func (s *RoleController) DeleteRole(c *gin.Context) {
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

	data, err := service.NewRoleService(s.svcCtx).DeleteRole(reqCtx, req)
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
// @Param		data	body		dto.RoleQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/role/find_role_list [POST]
func (s *RoleController) FindRoleList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RoleQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRoleService(s.svcCtx).FindRoleList(reqCtx, req)
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
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.RoleResourcesResp}	"返回信息"
// @Router		/admin_api/v1/role/find_role_resources [POST]
func (s *RoleController) FindRoleResources(c *gin.Context) {
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

	data, err := service.NewRoleService(s.svcCtx).FindRoleResources(reqCtx, req)
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
// @Param		data	body		dto.RoleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.RoleBackDTO}	"返回信息"
// @Router		/admin_api/v1/role/update_role [PUT]
func (s *RoleController) UpdateRole(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.RoleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRoleService(s.svcCtx).UpdateRole(reqCtx, req)
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
// @Param		data	body		dto.UpdateRoleApisReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/role/update_role_apis [POST]
func (s *RoleController) UpdateRoleApis(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UpdateRoleApisReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRoleService(s.svcCtx).UpdateRoleApis(reqCtx, req)
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
// @Param		data	body		dto.UpdateRoleMenusReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/role/update_role_menus [POST]
func (s *RoleController) UpdateRoleMenus(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.UpdateRoleMenusReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewRoleService(s.svcCtx).UpdateRoleMenus(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
