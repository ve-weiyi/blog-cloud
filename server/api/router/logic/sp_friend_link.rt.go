package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type FriendLinkRouter struct {
	svcCtx *svc.RouterContext
}

func NewFriendLinkRouter(svcCtx *svc.RouterContext) *FriendLinkRouter {
	return &FriendLinkRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 FriendLink 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *FriendLinkRouter) InitFriendLinkRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.FriendLinkController
	{
		loginRouter.POST("friend_link", handler.CreateFriendLink)                     // 新建FriendLink
		loginRouter.PUT("friend_link", handler.UpdateFriendLink)                      // 更新FriendLink
		loginRouter.DELETE("friend_link/:id", handler.DeleteFriendLink)               // 删除FriendLink
		loginRouter.DELETE("friend_link/batch_delete", handler.DeleteFriendLinkByIds) // 批量删除FriendLink列表

		publicRouter.GET("friend_link/:id", handler.FindFriendLink)       // 查询FriendLink
		publicRouter.POST("friend_link/list", handler.FindFriendLinkList) // 分页查询FriendLink列表
	}
}
