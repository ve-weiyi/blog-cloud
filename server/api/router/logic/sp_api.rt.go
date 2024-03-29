package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type ApiRouter struct {
	svcCtx *svc.RouterContext
}

func NewApiRouter(svcCtx *svc.RouterContext) *ApiRouter {
	return &ApiRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Api 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *ApiRouter) InitApiRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.ApiController
	{
		loginRouter.POST("api", handler.CreateApi)                     // 新建Api
		loginRouter.PUT("api", handler.UpdateApi)                      // 更新Api
		loginRouter.DELETE("api/:id", handler.DeleteApi)               // 删除Api
		loginRouter.DELETE("api/batch_delete", handler.DeleteApiByIds) // 批量删除Api列表

		loginRouter.GET("api/:id", handler.FindApi)                      // 查询Api
		loginRouter.POST("api/list", handler.FindApiList)                // 分页查询Api列表
		loginRouter.POST("api/sync", handler.SyncApiList)                // 同步Api列表
		loginRouter.POST("api/clean", handler.CreateApi)                 // 清空Api列表
		loginRouter.POST("api/details_list", handler.FindApiDetailsList) // 获取Api列表
	}
}
