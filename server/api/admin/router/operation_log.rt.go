package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewOperationLogController(s.svcCtx)
		// 批量删除操作记录
		group.DELETE("/operation_log/batch_delete_operation_log", handler.BatchDeleteOperationLog)
		// 删除操作记录
		group.DELETE("/operation_log/delete_operation_log", handler.DeleteOperationLog)
		// 分页获取操作记录列表
		group.POST("/operation_log/find_operation_log_list", handler.FindOperationLogList)
	}
}
