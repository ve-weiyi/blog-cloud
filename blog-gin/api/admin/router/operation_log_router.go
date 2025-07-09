package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type OperationLogRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewOperationLogRouter(svcCtx *svctx.ServiceContext) *OperationLogRouter {
	return &OperationLogRouter{
		svcCtx: svcCtx,
	}
}

func (s *OperationLogRouter) Register(r *gin.RouterGroup) {
	// OperationLog
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewOperationLogController(s.svcCtx)
		// 删除操作记录
		group.DELETE("/operation_log/deletes_operation_log", h.DeletesOperationLog)
		// 分页获取操作记录列表
		group.POST("/operation_log/find_operation_log_list", h.FindOperationLogList)
	}
}
