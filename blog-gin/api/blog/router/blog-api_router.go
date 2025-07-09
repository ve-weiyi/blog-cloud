package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type BlogApiRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewBlogApiRouter(svcCtx *svctx.ServiceContext) *BlogApiRouter {
	return &BlogApiRouter{
		svcCtx: svcCtx,
	}
}

func (s *BlogApiRouter) Register(r *gin.RouterGroup) {
	// BlogApi
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewBlogApiController(s.svcCtx)
		// ping
		group.GET("/ping", h.Ping)
	}
}
