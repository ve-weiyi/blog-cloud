package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
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
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewWebsiteController(s.svcCtx)
		// 获取博客前台首页信息
		group.GET("/blog", handler.GetBlogHomeInfo)
		// 获取关于我的信息
		group.GET("/blog/about_me", handler.GetAboutMe)
	}
}
