package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [TerminalToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)

		h := handler.NewAlbumController(s.svcCtx)
		// 获取相册列表
		group.POST("/album/find_album_list", h.FindAlbumList)
		// 获取相册下的照片列表
		group.POST("/album/find_photo_list", h.FindPhotoList)
		// 获取相册
		group.POST("/album/get_album", h.GetAlbum)
	}
}
