package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewRemarkController(s.svcCtx)
		// 分页获取留言列表
		group.POST("/remark/find_remark_list", h.FindRemarkList)
	}
	// Remark
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewRemarkController(s.svcCtx)
		// 创建留言
		group.POST("/remark/add_remark", h.AddRemark)
	}
}
