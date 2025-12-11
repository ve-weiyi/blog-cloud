package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type ApiRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewApiRouter(svcCtx *svctx.ServiceContext) *ApiRouter {
	return &ApiRouter{
		svcCtx: svcCtx,
	}
}

func (s *ApiRouter) Register(r *gin.RouterGroup) {
	// Api
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewApiController(s.svcCtx)
		// 创建api路由
		group.POST("/api/add_api", h.AddApi)
		// 清空接口列表
		group.POST("/api/clean_api_list", h.CleanApiList)
		// 删除api路由
		group.DELETE("/api/deletes_api", h.DeletesApi)
		// 分页获取api路由列表
		group.POST("/api/find_api_list", h.FindApiList)
		// 同步api列表
		group.POST("/api/sync_api_list", h.SyncApiList)
		// 更新api路由
		group.PUT("/api/update_api", h.UpdateApi)
	}
}
