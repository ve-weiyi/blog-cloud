package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
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

// 删除评论
func (s *CommentLogic) DeletesComment(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 查询评论列表(后台)
func (s *CommentLogic) FindCommentBackList(reqCtx *request.Context, in *types.QueryCommentReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 更新评论审核状态
func (s *CommentLogic) UpdateCommentStatus(reqCtx *request.Context, in *types.UpdateCommentStatusReq) (out *types.BatchResp, err error) {
	// todo

	return
}
