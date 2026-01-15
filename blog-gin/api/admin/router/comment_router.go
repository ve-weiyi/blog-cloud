package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/handler"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type CommentRouter struct {
	svcCtx *svctx.ServiceContext
}

func NewCommentRouter(svcCtx *svctx.ServiceContext) *CommentRouter {
	return &CommentRouter{
		svcCtx: svcCtx,
	}
}

func (s *CommentRouter) Register(r *gin.RouterGroup) {
	// Comment
	// [AdminToken Permission OperationLog]
	{
		group := r.Group("/admin-api/v1")
		group.Use(s.svcCtx.AdminToken)
		group.Use(s.svcCtx.Permission)
		group.Use(s.svcCtx.OperationLog)

		h := handler.NewCommentController(s.svcCtx)
		// 删除评论
		group.DELETE("/comment/deletes_comment", h.DeletesComment)
		// 查询评论列表(后台)
		group.POST("/comment/find_comment_back_list", h.FindCommentBackList)
		// 更新评论状态
		group.PUT("/comment/update_comment_status", h.UpdateCommentStatus)
	}
}
