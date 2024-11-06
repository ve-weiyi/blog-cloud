package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type FriendRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewFriendRouter(svcCtx *svctx.ServiceContext) *FriendRouter {
	return &FriendRouter{
		svcCtx: svcCtx,
	}
}

func (s *FriendRouter) Register(r *gin.RouterGroup) {
	// Friend
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewFriendController(s.svcCtx)
		// 分页获取友链列表
		group.POST("/friend_link/find_friend_list", handler.FindFriendList)
	}
}
