package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type AlbumRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewAlbumRouter(svcCtx *svctx.ServiceContext) *AlbumRouter {
	return &AlbumRouter{
		svcCtx: svcCtx,
	}
}

func (s *AlbumRouter) Register(r *gin.RouterGroup) {
	// Album
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewAlbumController(s.svcCtx)
		// 获取相册列表
		group.POST("/album/find_album_list", handler.FindAlbumList)
		// 获取相册下的照片列表
		group.POST("/album/find_photo_list", handler.FindPhotoList)
		// 获取相册
		group.POST("/album/get_album", handler.GetAlbum)
	}
}
