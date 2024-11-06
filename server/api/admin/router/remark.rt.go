package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewRemarkController(s.svcCtx)
		// 分页获取留言列表
		group.POST("/remark/find_remark_list", handler.FindRemarkList)
	}
	// Remark
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewRemarkController(s.svcCtx)
		// 批量删除留言
		group.DELETE("/remark/batch_delete_remark", handler.BatchDeleteRemark)
		// 删除留言
		group.DELETE("/remark/delete_remark", handler.DeleteRemark)
		// 更新留言
		group.PUT("/remark/update_remark", handler.UpdateRemark)
	}
}
