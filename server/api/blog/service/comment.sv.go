package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type CommentService struct {
	svcCtx *svctx.ServiceContext
}

func NewCommentService(svcCtx *svctx.ServiceContext) *CommentService {
	return &CommentService{
		svcCtx: svcCtx,
	}
}

// 查询评论列表
func (s *CommentService) FindCommentList(reqCtx *request.Context, in *dto.CommentQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询最新评论回复列表
func (s *CommentService) FindCommentRecentList(reqCtx *request.Context, in *dto.CommentQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 查询评论回复列表
func (s *CommentService) FindCommentReplyList(reqCtx *request.Context, in *dto.CommentQueryReq) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建评论
func (s *CommentService) AddComment(reqCtx *request.Context, in *dto.CommentNewReq) (out *dto.Comment, err error) {
	// todo

	return
}

// 点赞评论
func (s *CommentService) LikeComment(reqCtx *request.Context, in *dto.IdReq) (out *dto.EmptyResp, err error) {
	// todo

	return
}
