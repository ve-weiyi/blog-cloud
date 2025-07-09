package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type PhotoRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewPhotoRouter(svcCtx *svctx.ServiceContext) *PhotoRouter {
	return &PhotoRouter{
		svcCtx: svcCtx,
	}
}

func (s *PhotoRouter) Register(r *gin.RouterGroup) {
	// Photo
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewPhotoController(s.svcCtx)
		// 创建照片
		group.POST("/photo/add_photo", h.AddPhoto)
		// 删除照片
		group.DELETE("/photo/deletes_photo", h.DeletesPhoto)
		// 分页获取照片列表
		group.POST("/photo/find_photo_list", h.FindPhotoList)
		// 预删除照片
		group.PUT("/photo/pre_delete_photo", h.PreDeletePhoto)
		// 更新照片
		group.PUT("/photo/update_photo", h.UpdatePhoto)
	}
}
