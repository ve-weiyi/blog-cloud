package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateCommentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentLogic {
	return &UpdateCommentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论
func (l *UpdateCommentLogic) UpdateComment(in *discussionrpc.UpdateCommentRequest) (*discussionrpc.UpdateCommentResponse, error) {
	entity, err := l.svcCtx.TCommentModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.ParentId = in.ParentId
	entity.ReplyId = in.ReplyId
	entity.ReplyUserId = in.ReplyUserId
	entity.CommentContent = in.CommentContent
	entity.Type = in.Type
	entity.Status = in.Status

	_, err = l.svcCtx.TCommentModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.UpdateCommentResponse{Success: true}, nil
}
