package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UploadLogRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewUploadLogRouter(svcCtx *svctx.ServiceContext) *UploadLogRouter {
	return &UploadLogRouter{
		svcCtx: svcCtx,
	}
}

func (s *UploadLogRouter) Register(r *gin.RouterGroup) {
	// UploadLog
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewUploadLogController(s.svcCtx)
		// 删除登录日志
		group.DELETE("/upload_log/deletes_upload_log", h.DeletesUploadLog)
		// 查询登录日志
		group.POST("/user/find_upload_log_list", h.FindUploadLogList)
	}
}
