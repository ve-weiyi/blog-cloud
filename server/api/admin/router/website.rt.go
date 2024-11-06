package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type WebsiteRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewWebsiteRouter(svcCtx *svctx.ServiceContext) *WebsiteRouter {
	return &WebsiteRouter{
		svcCtx: svcCtx,
	}
}

func (s *WebsiteRouter) Register(r *gin.RouterGroup) {
	// Website
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewWebsiteController(s.svcCtx)
		// 获取后台首页信息
		group.GET("/admin", handler.GetAdminHomeInfo)
		// 获取关于我的信息
		group.GET("/admin/about_me", handler.GetAboutMe)
		// 更新关于我的信息
		group.PUT("/admin/about_me", handler.UpdateAboutMe)
		// 获取网站配置
		group.GET("/admin/get_website_config", handler.GetWebsiteConfig)
		// 获取服务器信息
		group.GET("/admin/system_state", handler.GetSystemState)
		// 更新网站配置
		group.PUT("/admin/update_website_config", handler.UpdateWebsiteConfig)
	}
}
