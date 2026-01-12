package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type FileLogRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewFileLogRouter(svcCtx *svctx.ServiceContext) *FileLogRouter {
	return &FileLogRouter{
		svcCtx: svcCtx,
	}
}

func (s *FileLogRouter) Register(r *gin.RouterGroup) {
	// FileLog
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewFileLogController(s.svcCtx)
		// 删除文件日志
		group.DELETE("/file_log/deletes_file_log", h.DeletesFileLog)
		// 查询文件日志
		group.POST("/file_log/find_file_log_list", h.FindFileLogList)
	}
}
