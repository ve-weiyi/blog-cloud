package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type AdminApiRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewAdminApiRouter(svcCtx *svctx.ServiceContext) *AdminApiRouter {
	return &AdminApiRouter{
		svcCtx: svcCtx,
	}
}

func (s *AdminApiRouter) Register(r *gin.RouterGroup) {
	// AdminApi
	// []
	{
		group := r.Group("/admin-api/v1")

		h := handler.NewAdminApiController(s.svcCtx)
		// ping
		group.GET("/ping", h.Ping)
	}
}
