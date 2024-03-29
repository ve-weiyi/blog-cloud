package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type PhotoRouter struct {
	svcCtx *svc.RouterContext
}

func NewPhotoRouter(svcCtx *svc.RouterContext) *PhotoRouter {
	return &PhotoRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Photo 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *PhotoRouter) InitPhotoRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.PhotoController
	{
		publicRouter.POST("photo", handler.CreatePhoto)       // 新建Photo
		publicRouter.PUT("photo", handler.UpdatePhoto)        // 更新Photo
		publicRouter.DELETE("photo/:id", handler.DeletePhoto) // 删除Photo
		publicRouter.GET("photo/:id", handler.FindPhoto)      // 查询Photo

		publicRouter.DELETE("photo/batch_delete", handler.DeletePhotoByIds) // 批量删除Photo列表
		publicRouter.POST("photo/list", handler.FindPhotoList)              // 分页查询Photo列表
	}
}
