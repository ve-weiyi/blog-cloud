package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type ArticleRouter struct {
	svcCtx *svc.RouterContext
}

func NewArticleRouter(svcCtx *svc.RouterContext) *ArticleRouter {
	return &ArticleRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Article 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *ArticleRouter) InitArticleRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.ArticleController
	// 后台操作
	{
		loginRouter.POST("/admin/article", handler.SaveArticle)          // 新建文章
		loginRouter.DELETE("/admin/article/:id", handler.DeleteArticle)  // 删除文章
		loginRouter.GET("/admin/article/:id", handler.FindArticle)       // 查询文章
		loginRouter.POST("/admin/article/list", handler.FindArticleList) // 分页查询文章列表

		loginRouter.PUT("/admin/article/top", handler.UpdateArticleTop)       // 置顶文章
		loginRouter.PUT("/admin/article/delete", handler.UpdateArticleDelete) // 逻辑删除文章(假删除)
	}
	// 前台操作接口
	{
		publicRouter.POST("article/archives", handler.FindArticleArchives)  // 文章归档
		publicRouter.POST("article/series", handler.FindArticleSeries)      // 根据条件获取Article列表
		publicRouter.POST("article/list", handler.FindArticleHomeList)      // 首页文章列表
		publicRouter.GET("article/:id/details", handler.FindArticleDetails) // 查询文章详情
		publicRouter.PUT("article/:id/like", handler.LikeArticle)           // 点赞文章
	}
}
