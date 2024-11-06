package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type ChatRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewChatRouter(svcCtx *svctx.ServiceContext) *ChatRouter {
	return &ChatRouter{
		svcCtx: svcCtx,
	}
}

func (s *ChatRouter) Register(r *gin.RouterGroup) {
	// Chat
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewChatController(s.svcCtx)
		// 查询聊天记录
		group.POST("/chat/messages", handler.GetChatMessages)
	}
}
