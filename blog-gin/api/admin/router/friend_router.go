package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
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
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewFriendController(s.svcCtx)
		// 创建友链
		group.POST("/friend/add_friend", h.AddFriend)
		// 删除友链
		group.DELETE("/friend/deletes_friend", h.DeletesFriend)
		// 分页获取友链列表
		group.POST("/friend/find_friend_list", h.FindFriendList)
		// 更新友链
		group.PUT("/friend/update_friend", h.UpdateFriend)
	}
}
