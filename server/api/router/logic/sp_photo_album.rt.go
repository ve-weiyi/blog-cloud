package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type PhotoAlbumRouter struct {
	svcCtx *svc.RouterContext
}

func NewPhotoAlbumRouter(svcCtx *svc.RouterContext) *PhotoAlbumRouter {
	return &PhotoAlbumRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 PhotoAlbum 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *PhotoAlbumRouter) InitPhotoAlbumRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.PhotoAlbumController
	{
		publicRouter.POST("photo_album", handler.CreatePhotoAlbum)                     // 新建PhotoAlbum
		publicRouter.PUT("photo_album", handler.UpdatePhotoAlbum)                      // 更新PhotoAlbum
		publicRouter.DELETE("photo_album/:id", handler.DeletePhotoAlbum)               // 删除PhotoAlbum
		publicRouter.DELETE("photo_album/batch_delete", handler.DeletePhotoAlbumByIds) // 批量删除PhotoAlbum列表

		publicRouter.GET("photo_album/:id", handler.FindPhotoAlbum)       // 查询PhotoAlbum
		publicRouter.POST("photo_album/list", handler.FindPhotoAlbumList) // 分页查询PhotoAlbum列表

		publicRouter.GET("photo_album/:id/details", handler.FindPhotoAlbumDetails)       // 获取PhotoAlbum详情
		publicRouter.POST("photo_album/details_list", handler.FindPhotoAlbumDetailsList) // 获取PhotoAlbum详情列表
	}
}
