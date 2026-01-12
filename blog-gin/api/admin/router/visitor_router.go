package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type VisitorRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewVisitorRouter(svcCtx *svctx.ServiceContext) *VisitorRouter {
	return &VisitorRouter{
		svcCtx: svcCtx,
	}
}

func (s *VisitorRouter) Register(r *gin.RouterGroup) {
	// Visitor
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewVisitorController(s.svcCtx)
		// 分页获取游客列表
		group.POST("/visitor/find_visitor_list", h.FindVisitorList)
	}
}
