package admin

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/router"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

// register your handlers here
func RegisterHandlers(r *gin.RouterGroup, svCtx *svctx.ServiceContext) {
	router.NewAccountRouter(svCtx).Register(r)
	router.NewAlbumRouter(svCtx).Register(r)
	router.NewApiRouter(svCtx).Register(r)
	router.NewArticleRouter(svCtx).Register(r)
	router.NewAuthRouter(svCtx).Register(r)
	router.NewBannerRouter(svCtx).Register(r)
	router.NewCategoryRouter(svCtx).Register(r)
	router.NewCommentRouter(svCtx).Register(r)
	router.NewCommonRouter(svCtx).Register(r)
	router.NewFileRouter(svCtx).Register(r)
	router.NewFriendRouter(svCtx).Register(r)
	router.NewMenuRouter(svCtx).Register(r)
	router.NewOperationLogRouter(svCtx).Register(r)
	router.NewPhotoRouter(svCtx).Register(r)
	router.NewRemarkRouter(svCtx).Register(r)
	router.NewRoleRouter(svCtx).Register(r)
	router.NewTagRouter(svCtx).Register(r)
	router.NewTalkRouter(svCtx).Register(r)
	router.NewUserRouter(svCtx).Register(r)
	router.NewWebsiteRouter(svCtx).Register(r)
}
