package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type UploadRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewUploadRouter(svcCtx *svctx.ServiceContext) *UploadRouter {
	return &UploadRouter{
		svcCtx: svcCtx,
	}
}

func (s *UploadRouter) Register(r *gin.RouterGroup) {
	// Upload
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewUploadController(s.svcCtx)
		// 删除文件列表
		group.DELETE("/upload/deletes_upload_file", h.DeletesUploadFile)
		// 获取文件列表
		group.POST("/upload/list_upload_file", h.ListUploadFile)
		// 上传文件列表
		group.POST("/upload/multi_upload_file", h.MultiUploadFile)
		// 上传文件
		group.POST("/upload/upload_file", h.UploadFile)
	}
}
