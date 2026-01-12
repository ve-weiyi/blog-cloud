package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
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
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewPageController(s.svcCtx)
		// 分页获取页面列表
		group.POST("/page/find_page_list", h.FindPageList)
	}
}
