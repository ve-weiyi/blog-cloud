package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewPhotoController(s.svcCtx)
		// 分页获取照片列表
		group.POST("/photo/find_photo_list", handler.FindPhotoList)
	}
	// Photo
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewPhotoController(s.svcCtx)
		// 批量删除照片
		group.DELETE("/album/batch_delete_photo", handler.BatchDeletePhoto)
		// 创建照片
		group.POST("/photo/add_photo", handler.AddPhoto)
		// 删除照片
		group.DELETE("/photo/delete_photo", handler.DeletePhoto)
		// 更新照片
		group.PUT("/photo/update_photo", handler.UpdatePhoto)
	}
}
