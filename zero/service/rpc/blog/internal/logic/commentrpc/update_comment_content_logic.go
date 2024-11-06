package commentrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCommentContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCommentContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCommentContentLogic {
	return &UpdateCommentContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新评论
func (l *UpdateCommentContentLogic) UpdateCommentContent(in *commentrpc.UpdateCommentContentReq) (*commentrpc.CommentDetails, error) {
	entity, err := l.svcCtx.TCommentModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.CommentContent = in.CommentContent

	_, err = l.svcCtx.TCommentModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertCommentOut(entity), nil
}
