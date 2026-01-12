package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewRemarkController(s.svcCtx)
		// 删除留言
		group.DELETE("/remark/deletes_remark", h.DeletesRemark)
		// 分页获取留言列表
		group.POST("/remark/find_remark_list", h.FindRemarkList)
		// 更新留言状态
		group.PUT("/remark/update_remark_status", h.UpdateRemarkStatus)
	}
}
