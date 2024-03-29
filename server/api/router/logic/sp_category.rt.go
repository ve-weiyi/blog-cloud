package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type CategoryRouter struct {
	svcCtx *svc.RouterContext
}

func NewCategoryRouter(svcCtx *svc.RouterContext) *CategoryRouter {
	return &CategoryRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Category 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *CategoryRouter) InitCategoryRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.CategoryController
	{
		loginRouter.POST("category", handler.CreateCategory)                     // 新建Category
		loginRouter.PUT("category", handler.UpdateCategory)                      // 更新Category
		loginRouter.DELETE("category/:id", handler.DeleteCategory)               // 删除Category
		loginRouter.DELETE("category/batch_delete", handler.DeleteCategoryByIds) // 批量删除Category列表

		publicRouter.GET("category/:id", handler.FindCategory)                      // 查询Category
		publicRouter.POST("category/list", handler.FindCategoryList)                // 分页查询Category列表
		publicRouter.POST("category/details_list", handler.FindCategoryDetailsList) // 查询Category详情列表
	}
}
