package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type WebsocketRouter struct {
	svcCtx *svc.RouterContext
}

func NewWebsocketRouter(ctx *svc.RouterContext) *WebsocketRouter {
	return &WebsocketRouter{
		svcCtx: ctx,
	}
}

// 初始化 Blog 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *WebsocketRouter) InitWebsocketRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.WebsocketController
	{
		publicRouter.GET("ws", handler.WebSocket) // websocket
	}
}
