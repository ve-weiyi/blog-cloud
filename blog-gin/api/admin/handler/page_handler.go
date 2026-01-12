package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type PageController struct {
	svcCtx *svctx.ServiceContext
}

func NewPageController(svcCtx *svctx.ServiceContext) *PageController {
	return &PageController{
		svcCtx: svcCtx,
	}
}

// @Tags		Page
// @Summary		"创建页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewPageReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageBackVO}	"返回信息"
// @Router		/admin-api/v1/page/add_page [POST]
func (s *PageController) AddPage(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewPageReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPageLogic(s.svcCtx).AddPage(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Page
// @Summary		"删除页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/page/delete_page [DELETE]
func (s *PageController) DeletePage(c *gin.Context) {
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

	data, err := logic.NewPageLogic(s.svcCtx).DeletePage(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Page
// @Summary		"分页获取页面列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.QueryPageReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/page/find_page_list [POST]
func (s *PageController) FindPageList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryPageReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPageLogic(s.svcCtx).FindPageList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Page
// @Summary		"更新页面"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.NewPageReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageBackVO}	"返回信息"
// @Router		/admin-api/v1/page/update_page [PUT]
func (s *PageController) UpdatePage(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewPageReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewPageLogic(s.svcCtx).UpdatePage(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
