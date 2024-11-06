package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewCategoryController(s.svcCtx)
		// 分页获取文章分类列表
		group.POST("/category/find_category_list", handler.FindCategoryList)
	}
	// Category
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewCategoryController(s.svcCtx)
		// 创建文章分类
		group.POST("/category/add_category", handler.AddCategory)
		// 批量删除文章分类
		group.DELETE("/category/batch_delete_category", handler.BatchDeleteCategory)
		// 删除文章分类
		group.DELETE("/category/delete_category", handler.DeleteCategory)
		// 更新文章分类
		group.PUT("/category/update_category", handler.UpdateCategory)
	}
}
