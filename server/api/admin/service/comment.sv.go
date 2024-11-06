package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
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

// 批量删除评论
func (s *CommentService) BatchDeleteComment(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除评论
func (s *CommentService) DeleteComment(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 查询评论列表(后台)
func (s *CommentService) FindCommentBackList(reqCtx *request.Context, in *dto.CommentQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 更新评论审核状态
func (s *CommentService) UpdateCommentReview(reqCtx *request.Context, in *dto.CommentReviewReq) (out *dto.BatchResp, err error) {
	// todo

	return
}
