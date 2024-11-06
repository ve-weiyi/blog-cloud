package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TagRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewTagRouter(svcCtx *svctx.ServiceContext) *TagRouter {
	return &TagRouter{
		svcCtx: svcCtx,
	}
}

func (s *TagRouter) Register(r *gin.RouterGroup) {
	// Tag
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewTagController(s.svcCtx)
		// 分页获取标签列表
		group.POST("/tag/find_tag_list", handler.FindTagList)
	}
}
