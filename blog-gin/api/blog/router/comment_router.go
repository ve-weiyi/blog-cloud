package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/handler"
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
	// [TerminalToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)

		h := handler.NewCommentController(s.svcCtx)
		// 查询评论列表
		group.POST("/comment/find_comment_list", h.FindCommentList)
		// 查询最新评论回复列表
		group.POST("/comment/find_comment_recent_list", h.FindCommentRecentList)
		// 查询评论回复列表
		group.POST("/comment/find_comment_reply_list", h.FindCommentReplyList)
	}
	// Comment
	// [TerminalToken UserToken]
	{
		group := r.Group("/blog-api/v1")
		group.Use(s.svcCtx.TerminalToken)
		group.Use(s.svcCtx.UserToken)

		h := handler.NewCommentController(s.svcCtx)
		// 创建评论
		group.POST("/comment/add_comment", h.AddComment)
		// 点赞评论
		group.POST("/comment/like_comment", h.LikeComment)
		// 更新评论
		group.POST("/comment/update_comment", h.UpdateComment)
	}
}
