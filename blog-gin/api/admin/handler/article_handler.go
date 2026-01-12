package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/logic"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/response"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
// @Param		data	body		types.NewArticleReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.ArticleBackVO}	"返回信息"
// @Router		/admin-api/v1/article/add_article [POST]
func (s *ArticleController) AddArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewArticleReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewArticleLogic(s.svcCtx).AddArticle(reqCtx, req)
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
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.BatchResp}	"返回信息"
// @Router		/admin-api/v1/article/delete_article [DELETE]
func (s *ArticleController) DeleteArticle(c *gin.Context) {
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

	data, err := logic.NewArticleLogic(s.svcCtx).DeleteArticle(reqCtx, req)
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
// @Param		data	body		types.IdsReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/article/export_article_list [POST]
func (s *ArticleController) ExportArticleList(c *gin.Context) {
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

	data, err := logic.NewArticleLogic(s.svcCtx).ExportArticleList(reqCtx, req)
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
// @Param		data	body		types.QueryArticleReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.PageResp}	"返回信息"
// @Router		/admin-api/v1/article/find_article_list [POST]
func (s *ArticleController) FindArticleList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.QueryArticleReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewArticleLogic(s.svcCtx).FindArticleList(reqCtx, req)
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
// @Param		data	body		types.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.ArticleBackVO}	"返回信息"
// @Router		/admin-api/v1/article/get_article [POST]
func (s *ArticleController) GetArticle(c *gin.Context) {
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

	data, err := logic.NewArticleLogic(s.svcCtx).GetArticle(reqCtx, req)
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
// @Param		data	body		types.NewArticleReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.ArticleBackVO}	"返回信息"
// @Router		/admin-api/v1/article/update_article [PUT]
func (s *ArticleController) UpdateArticle(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.NewArticleReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewArticleLogic(s.svcCtx).UpdateArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"更新文章删除状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateArticleDeleteReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/article/update_article_delete [PUT]
func (s *ArticleController) UpdateArticleDelete(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateArticleDeleteReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewArticleLogic(s.svcCtx).UpdateArticleDelete(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"更新文章置顶状态"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		types.UpdateArticleTopReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=types.EmptyResp}	"返回信息"
// @Router		/admin-api/v1/article/update_article_top [PUT]
func (s *ArticleController) UpdateArticleTop(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *types.UpdateArticleTopReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := logic.NewArticleLogic(s.svcCtx).UpdateArticleTop(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	response.ResponseOk(c, data)
}
