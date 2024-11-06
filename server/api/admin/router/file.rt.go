package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewFileController(s.svcCtx)
		// 分页获取文件列表
		group.POST("/file/find_file_list", handler.FindFileList)
	}
	// File
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewFileController(s.svcCtx)
		// 创建文件目录
		group.POST("/file/add_file_folder", handler.AddFileFolder)
		// 删除文件列表
		group.DELETE("/file/deletes_file", handler.DeletesFile)
		// 上传文件列表
		group.POST("/file/multi_upload_file", handler.MultiUploadFile)
		// 上传文件
		group.POST("/file/upload_file", handler.UploadFile)
	}
}
