package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type RemarkRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewRemarkRouter(svcCtx *svctx.ServiceContext) *RemarkRouter {
	return &RemarkRouter{
		svcCtx: svcCtx,
	}
}

func (s *RemarkRouter) Register(r *gin.RouterGroup) {
	// Remark
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewRemarkController(s.svcCtx)
		// 分页获取留言列表
		group.POST("/remark/find_remark_list", handler.FindRemarkList)
	}
	// Remark
	// [SignToken JwtToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewRemarkController(s.svcCtx)
		// 创建留言
		group.POST("/remark/add_remark", handler.AddRemark)
	}
}
