package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewBannerController(s.svcCtx)
		// 创建页面
		group.POST("/banner/add_banner", handler.AddBanner)
		// 删除页面
		group.DELETE("/banner/delete_banner", handler.DeleteBanner)
		// 分页获取页面列表
		group.POST("/banner/find_banner_list", handler.FindBannerList)
		// 更新页面
		group.PUT("/banner/update_banner", handler.UpdateBanner)
	}
}
