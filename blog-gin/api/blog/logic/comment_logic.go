package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/blog/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type CommentLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewCommentLogic(svcCtx *svctx.ServiceContext) *CommentLogic {
	return &CommentLogic{
		svcCtx: svcCtx,
	}
}

// 查询评论列表
func (s *CommentLogic) FindCommentList(reqCtx *request.Context, in *types.QueryCommentReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询最新评论回复列表
func (s *CommentLogic) FindCommentRecentList(reqCtx *request.Context, in *types.QueryCommentReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 查询评论回复列表
func (s *CommentLogic) FindCommentReplyList(reqCtx *request.Context, in *types.QueryCommentReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 创建评论
func (s *CommentLogic) AddComment(reqCtx *request.Context, in *types.NewCommentReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 点赞评论
func (s *CommentLogic) LikeComment(reqCtx *request.Context, in *types.IdReq) (out *types.EmptyResp, err error) {
	// todo

	return
}

// 更新评论
func (s *CommentLogic) UpdateComment(reqCtx *request.Context, in *types.UpdateCommentReq) (out *types.EmptyResp, err error) {
	// todo

	return
}
