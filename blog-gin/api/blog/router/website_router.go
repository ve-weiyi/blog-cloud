package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
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
	// [TerminalToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)

		h := handler.NewWebsiteController(s.svcCtx)
		// 获取博客前台首页信息
		group.GET("/blog", h.GetBlogHomeInfo)
		// 获取关于我的信息
		group.GET("/blog/about_me", h.GetAboutMe)
	}
}
