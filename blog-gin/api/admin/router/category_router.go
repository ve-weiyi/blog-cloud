package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewCategoryController(s.svcCtx)
		// 创建文章分类
		group.POST("/category/add_category", h.AddCategory)
		// 删除文章分类
		group.DELETE("/category/deletes_category", h.DeletesCategory)
		// 分页获取文章分类列表
		group.POST("/category/find_category_list", h.FindCategoryList)
		// 更新文章分类
		group.PUT("/category/update_category", h.UpdateCategory)
	}
}
