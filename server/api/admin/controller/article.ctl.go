package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/service"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/response"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type ArticleController struct {
	svcCtx *svctx.ServiceContext
}

func NewArticleController(svcCtx *svctx.ServiceContext) *ArticleController {
	return &ArticleController{
		svcCtx: svcCtx,
	}
}

// @Tags		Article
// @Summary		"添加文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.ArticleBackDTO}	"返回信息"
// @Router		/admin_api/v1/article/add_article [POST]
func (s *ArticleController) AddArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).AddArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"删除文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.BatchResp}	"返回信息"
// @Router		/admin_api/v1/article/delete_article [POST]
func (s *ArticleController) DeleteArticle(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).DeleteArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"导出文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/article/export_article_list [POST]
func (s *ArticleController) ExportArticleList(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).ExportArticleList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"查询文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleQuery		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/admin_api/v1/article/find_article_list [POST]
func (s *ArticleController) FindArticleList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleQuery
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).FindArticleList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"查询文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.ArticleBackDTO}	"返回信息"
// @Router		/admin_api/v1/article/get_article [POST]
func (s *ArticleController) GetArticle(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).GetArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"回收文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleRecycleReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/article/recycle_article [POST]
func (s *ArticleController) RecycleArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleRecycleReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).RecycleArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"置顶文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleTopReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/admin_api/v1/article/top_article [POST]
func (s *ArticleController) TopArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleTopReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).TopArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"保存文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleNewReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.ArticleBackDTO}	"返回信息"
// @Router		/admin_api/v1/article/update_article [POST]
func (s *ArticleController) UpdateArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleNewReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).UpdateArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
