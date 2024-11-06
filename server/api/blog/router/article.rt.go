package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type ArticleRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewArticleRouter(svcCtx *svctx.ServiceContext) *ArticleRouter {
	return &ArticleRouter{
		svcCtx: svcCtx,
	}
}

func (s *ArticleRouter) Register(r *gin.RouterGroup) {
	// Article
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewArticleController(s.svcCtx)
		// 文章归档(时间轴)
		group.POST("/article/get_article_archives", handler.FindArticleArchives)
		// 通过分类获取文章列表
		group.POST("/article/get_article_classify_category", handler.FindArticleClassifyCategory)
		// 通过标签获取文章列表
		group.POST("/article/get_article_classify_tag", handler.FindArticleClassifyTag)
		// 获取文章详情
		group.POST("/article/get_article_details", handler.GetArticleDetails)
		// 获取首页文章列表
		group.POST("/article/get_article_home_list", handler.FindArticleHomeList)
		// 获取首页推荐文章列表
		group.POST("/article/get_article_recommend", handler.FindArticleRecommend)
	}
	// Article
	// [SignToken JwtToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewArticleController(s.svcCtx)
		// 点赞文章
		group.POST("/article/like_article", handler.LikeArticle)
	}
}
