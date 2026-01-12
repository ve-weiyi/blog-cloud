package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewWebsiteController(s.svcCtx)
		// 获取后台首页信息
		group.GET("/admin", h.GetAdminHomeInfo)
		// 获取关于我的信息
		group.GET("/admin/get_about_me", h.GetAboutMe)
		// 获取服务器信息
		group.GET("/admin/get_system_state", h.GetSystemState)
		// 获取用户分布地区
		group.POST("/admin/get_user_area_stats", h.GetUserAreaStats)
		// 获取访客数据分析
		group.GET("/admin/get_visit_stats", h.GetVisitStats)
		// 获取访客数据趋势
		group.POST("/admin/get_visit_trend", h.GetVisitTrend)
		// 获取网站配置
		group.GET("/admin/get_website_config", h.GetWebsiteConfig)
		// 更新关于我的信息
		group.PUT("/admin/update_about_me", h.UpdateAboutMe)
		// 更新网站配置
		group.PUT("/admin/update_website_config", h.UpdateWebsiteConfig)
	}
}
