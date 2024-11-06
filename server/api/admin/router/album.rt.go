package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewAlbumController(s.svcCtx)
		// 分页获取相册列表
		group.POST("/album/find_album_list", handler.FindAlbumList)
	}
	// Album
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewAlbumController(s.svcCtx)
		// 创建相册
		group.POST("/album/add_album", handler.AddAlbum)
		// 批量删除相册
		group.DELETE("/album/batch_delete_album", handler.BatchDeleteAlbum)
		// 删除相册
		group.DELETE("/album/delete_album", handler.DeleteAlbum)
		// 查询相册
		group.POST("/album/get_album", handler.GetAlbum)
		// 更新相册
		group.PUT("/album/update_album", handler.UpdateAlbum)
	}
}
