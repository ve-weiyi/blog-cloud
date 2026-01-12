package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewArticleController(s.svcCtx)
		// 文章归档(时间轴)
		group.POST("/article/find_article_archives", h.FindArticleArchives)
		// 通过分类获取文章列表
		group.POST("/article/find_article_classify_category", h.FindArticleClassifyCategory)
		// 通过标签获取文章列表
		group.POST("/article/find_article_classify_tag", h.FindArticleClassifyTag)
		// 获取首页文章列表
		group.POST("/article/find_article_home_list", h.FindArticleHomeList)
		// 获取首页推荐文章列表
		group.POST("/article/find_article_recommend", h.FindArticleRecommend)
		// 获取文章详情
		group.POST("/article/get_article_details", h.GetArticleDetails)
	}
	// Article
	// [UserToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.UserToken)

		h := handler.NewArticleController(s.svcCtx)
		// 点赞文章
		group.PUT("/article/like_article", h.LikeArticle)
	}
}
