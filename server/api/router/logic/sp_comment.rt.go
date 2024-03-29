package logic

import (
	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/server/api/router/svc"
)

type CommentRouter struct {
	svcCtx *svc.RouterContext
}

func NewCommentRouter(svcCtx *svc.RouterContext) *CommentRouter {
	return &CommentRouter{
		svcCtx: svcCtx,
	}
}

// 初始化 Comment 路由信息
// publicRouter 公开路由，不登录就可以访问
// loginRouter  登录路由，登录后才可以访问
func (s *CommentRouter) InitCommentRouter(publicRouter *gin.RouterGroup, loginRouter *gin.RouterGroup) {

	var handler = s.svcCtx.CommentController
	{
		loginRouter.POST("comment", handler.CreateComment)                     // 新建Comment
		loginRouter.PUT("comment", handler.UpdateComment)                      // 更新Comment
		loginRouter.DELETE("comment/:id", handler.DeleteComment)               // 删除Comment
		loginRouter.POST("comment/:id/like", handler.LikeComment)              // 点赞评论
		loginRouter.DELETE("comment/batch_delete", handler.DeleteCommentByIds) // 批量删除Comment列表

		loginRouter.GET("comment/:id", handler.FindComment)                // 查询Comment
		loginRouter.POST("comment/list/back", handler.FindCommentBackList) // 分页查询Comment列表
	}
	{
		publicRouter.POST("comment/list", handler.FindCommentList)                // 分页查询Comment列表
		publicRouter.POST("comment/details_list", handler.FindCommentDetailsList) // 分页获取评论列表
		publicRouter.POST("comment/:id/reply_list", handler.FindCommentReplyList) // 查询评论回复列表
	}
}
