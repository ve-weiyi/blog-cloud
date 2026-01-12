package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type ApiController struct {
	svcCtx *svctx.ServiceContext
}

func NewApiController(svcCtx *svctx.ServiceContext) *ApiController {
	return &ApiController{
		svcCtx: svcCtx,
	}
}

// @Tags		Api
// @Summary		"创建api路由"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewApiReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.ApiBackVO}	"返回信息"
// @Router		/admin-api/v1/api/add_api [POST]
func (s *ApiController) AddApi(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewApiReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewApiLogic(s.svcCtx).AddApi(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Api
// @Summary		"清空接口列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/api/clean_api_list [POST]
func (s *ApiController) CleanApiList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.EmptyReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewApiLogic(s.svcCtx).CleanApiList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Api
// @Summary		"删除api路由"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/api/deletes_api [DELETE]
func (s *ApiController) DeletesApi(c *gin.Context) {
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

	data, err := logic.NewApiLogic(s.svcCtx).DeletesApi(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Api
// @Summary		"分页获取api路由列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryApiReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/api/find_api_list [POST]
func (s *ApiController) FindApiList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryApiReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewApiLogic(s.svcCtx).FindApiList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Api
// @Summary		"同步api列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.SyncApiReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/api/sync_api_list [POST]
func (s *ApiController) SyncApiList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.SyncApiReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewApiLogic(s.svcCtx).SyncApiList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Api
// @Summary		"更新api路由"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewApiReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.ApiBackVO}	"返回信息"
// @Router		/admin-api/v1/api/update_api [PUT]
func (s *ApiController) UpdateApi(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewApiReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewApiLogic(s.svcCtx).UpdateApi(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
