package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
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
	// []
	{
		group := r.Group("/blog-api/v1")

		h := handler.NewFriendController(s.svcCtx)
		// 分页获取友链列表
		group.POST("/friend/find_friend_list", h.FindFriendList)
	}
}
