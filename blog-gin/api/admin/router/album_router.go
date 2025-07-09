package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewAlbumController(s.svcCtx)
		// 创建相册
		group.POST("/album/add_album", h.AddAlbum)
		// 删除相册
		group.DELETE("/album/deletes_album", h.DeletesAlbum)
		// 分页获取相册列表
		group.POST("/album/find_album_list", h.FindAlbumList)
		// 查询相册
		group.POST("/album/get_album", h.GetAlbum)
		// 预删除相册
		group.POST("/album/pre_delete_album", h.PreDeleteAlbum)
		// 更新相册
		group.PUT("/album/update_album", h.UpdateAlbum)
	}
}
