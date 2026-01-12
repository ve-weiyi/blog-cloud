package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewArticleController(s.svcCtx)
		// 添加文章
		group.POST("/article/add_article", h.AddArticle)
		// 删除文章
		group.DELETE("/article/delete_article", h.DeleteArticle)
		// 导出文章列表
		group.POST("/article/export_article_list", h.ExportArticleList)
		// 查询文章列表
		group.POST("/article/find_article_list", h.FindArticleList)
		// 查询文章
		group.POST("/article/get_article", h.GetArticle)
		// 保存文章
		group.PUT("/article/update_article", h.UpdateArticle)
		// 更新文章删除状态
		group.PUT("/article/update_article_delete", h.UpdateArticleDelete)
		// 更新文章置顶状态
		group.PUT("/article/update_article_top", h.UpdateArticleTop)
	}
}
