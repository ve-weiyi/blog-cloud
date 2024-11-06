package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
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
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewFriendController(s.svcCtx)
		// 分页获取友链列表
		group.POST("/friend/find_friend_list", handler.FindFriendList)
	}
	// Friend
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewFriendController(s.svcCtx)
		// 创建友链
		group.POST("/friend/add_friend", handler.AddFriend)
		// 批量删除友链
		group.DELETE("/friend/batch_delete_friend", handler.BatchDeleteFriend)
		// 删除友链
		group.DELETE("/friend/delete_friend", handler.DeleteFriend)
		// 更新友链
		group.PUT("/friend/update_friend", handler.UpdateFriend)
	}
}
