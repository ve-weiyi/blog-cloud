package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type CategoryRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewCategoryRouter(svcCtx *svctx.ServiceContext) *CategoryRouter {
	return &CategoryRouter{
		svcCtx: svcCtx,
	}
}

func (s *CategoryRouter) Register(r *gin.RouterGroup) {
	// Category
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewCategoryController(s.svcCtx)
		// 分页获取文章分类列表
		group.POST("/category/find_category_list", handler.FindCategoryList)
	}
}
