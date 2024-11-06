package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type BannerRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewBannerRouter(svcCtx *svctx.ServiceContext) *BannerRouter {
	return &BannerRouter{
		svcCtx: svcCtx,
	}
}

func (s *BannerRouter) Register(r *gin.RouterGroup) {
	// Banner
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewBannerController(s.svcCtx)
		// 分页获取页面列表
		group.POST("/banner/find_banner_list", handler.FindBannerList)
	}
}
