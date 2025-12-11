package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type PageRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewPageRouter(svcCtx *svctx.ServiceContext) *PageRouter {
	return &PageRouter{
		svcCtx: svcCtx,
	}
}

func (s *PageRouter) Register(r *gin.RouterGroup) {
	// Page
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewPageController(s.svcCtx)
		// 创建页面
		group.POST("/page/add_page", h.AddPage)
		// 删除页面
		group.DELETE("/page/delete_page", h.DeletePage)
		// 分页获取页面列表
		group.POST("/page/find_page_list", h.FindPageList)
		// 更新页面
		group.PUT("/page/update_page", h.UpdatePage)
	}
}
