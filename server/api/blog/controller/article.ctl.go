package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/service"
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
// @Summary		"文章归档(时间轴)"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleArchivesQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/article/get_article_archives [POST]
func (s *ArticleController) FindArticleArchives(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleArchivesQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).FindArticleArchives(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"通过分类获取文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleClassifyQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/article/get_article_classify_category [POST]
func (s *ArticleController) FindArticleClassifyCategory(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleClassifyQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).FindArticleClassifyCategory(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"通过标签获取文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleClassifyQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/article/get_article_classify_tag [POST]
func (s *ArticleController) FindArticleClassifyTag(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleClassifyQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).FindArticleClassifyTag(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"获取文章详情"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.ArticleDeatils}	"返回信息"
// @Router		/api/v1/article/get_article_details [POST]
func (s *ArticleController) GetArticleDetails(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).GetArticleDetails(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"获取首页文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.ArticleHomeQueryReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/article/get_article_home_list [POST]
func (s *ArticleController) FindArticleHomeList(c *gin.Context) {
	reqCtx, err := request.ParseRequestContext(c)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	var req *dto.ArticleHomeQueryReq
	err = request.ShouldBind(c, &req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	data, err := service.NewArticleService(s.svcCtx).FindArticleHomeList(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"获取首页推荐文章列表"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.EmptyReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.PageResp}	"返回信息"
// @Router		/api/v1/article/get_article_recommend [POST]
func (s *ArticleController) FindArticleRecommend(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).FindArticleRecommend(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}

// @Tags		Article
// @Summary		"点赞文章"
// @accept		application/json
// @Produce		application/json
// @Param		data	body		dto.IdReq		true	"请求参数"
// @Success		200		{object}	response.Body{data=dto.EmptyResp}	"返回信息"
// @Router		/api/v1/article/like_article [POST]
func (s *ArticleController) LikeArticle(c *gin.Context) {
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

	data, err := service.NewArticleService(s.svcCtx).LikeArticle(reqCtx, req)
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	response.ResponseOk(c, data)
}
