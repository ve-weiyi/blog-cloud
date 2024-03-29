package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type MenuRouter struct {
	svcCtx *svc.RouterContext
}

func NewMenuRouter(svcCtx *svc.RouterContext) *MenuRouter {
	return &MenuRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Menu 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *MenuRouter) InitMenuRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.MenuController
	{
		loginRouter.POST("menu", handler.CreateMenu)                     // 新建Menu
		loginRouter.PUT("menu", handler.UpdateMenu)                      // 更新Menu
		loginRouter.DELETE("menu/:id", handler.DeleteMenu)               // 删除Menu
		loginRouter.DELETE("menu/batch_delete", handler.DeleteMenuByIds) // 批量删除Menu列表

		loginRouter.GET("menu/:id", handler.FindMenu)       // 查询Menu
		loginRouter.POST("menu/list", handler.FindMenuList) // 分页查询Menu列表

		loginRouter.POST("menu/details_list", handler.FindMenuDetailsList)

		loginRouter.POST("menu/sync", handler.SyncMenuList)   // 同步Menu列表
		loginRouter.POST("menu/clean", handler.CleanMenuList) // 清空Menu列表
	}
}
