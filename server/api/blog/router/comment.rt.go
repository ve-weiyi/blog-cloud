package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/controller"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
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
	// [SignToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)

		handler := controller.NewCommentController(s.svcCtx)
		// 查询评论列表
		group.POST("/comment/find_comment_list", handler.FindCommentList)
		// 查询最新评论回复列表
		group.POST("/comment/find_comment_recent_list", handler.FindCommentRecentList)
		// 查询评论回复列表
		group.POST("/comment/find_comment_reply_list", handler.FindCommentReplyList)
	}
	// Comment
	// [SignToken JwtToken]
	{
		group := r.Group("/api/v1")
		group.Use(s.svcCtx.MiddlewareSignToken)
		group.Use(s.svcCtx.MiddlewareJwtToken)

		handler := controller.NewCommentController(s.svcCtx)
		// 创建评论
		group.POST("/comment/add_comment", handler.AddComment)
		// 点赞评论
		group.POST("/comment/like_comment", handler.LikeComment)
	}
}
