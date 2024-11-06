package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewArticleController(s.svcCtx)
		// 添加文章
		group.POST("/article/add_article", handler.AddArticle)
		// 删除文章
		group.POST("/article/delete_article", handler.DeleteArticle)
		// 导出文章列表
		group.POST("/article/export_article_list", handler.ExportArticleList)
		// 查询文章列表
		group.POST("/article/find_article_list", handler.FindArticleList)
		// 查询文章
		group.POST("/article/get_article", handler.GetArticle)
		// 回收文章
		group.POST("/article/recycle_article", handler.RecycleArticle)
		// 置顶文章
		group.POST("/article/top_article", handler.TopArticle)
		// 保存文章
		group.POST("/article/update_article", handler.UpdateArticle)
	}
}
