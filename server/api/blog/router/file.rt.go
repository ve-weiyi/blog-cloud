package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type FileRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewFileRouter(svcCtx *svctx.ServiceContext) *FileRouter {
	return &FileRouter{
		svcCtx: svcCtx,
	}
}

func (s *FileRouter) Register(r *gin.RouterGroup) {
	// File
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewFileController(s.svcCtx)
		// 上传文件列表
		group.POST("/file/multi_upload_file", handler.MultiUploadFile)
		// 上传文件
		group.POST("/file/upload_file", handler.UploadFile)
	}
}
