package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type TagRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewTagRouter(svcCtx *svctx.ServiceContext) *TagRouter {
	return &TagRouter{
		svcCtx: svcCtx,
	}
}

func (s *TagRouter) Register(r *gin.RouterGroup) {
	// Tag
	// [SignToken]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewTagController(s.svcCtx)
		// 分页获取标签列表
		group.POST("/tag/find_tag_list", handler.FindTagList)
	}
	// Tag
	// [SignToken JwtToken Operation]
	{
		group := r.Group("/admin_api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)
		group.Use(s.svcCtx.MiddlewareOperation)

		handler := controller.NewTagController(s.svcCtx)
		// 创建标签
		group.POST("/tag/add_tag", handler.AddTag)
		// 批量删除标签
		group.DELETE("/tag/batch_delete_tag", handler.BatchDeleteTag)
		// 删除标签
		group.DELETE("/tag/delete_tag", handler.DeleteTag)
		// 更新标签
		group.PUT("/tag/update_tag", handler.UpdateTag)
	}
}
