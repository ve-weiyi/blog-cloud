package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteCommentsLogic {
	return &BatchDeleteCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除评论
func (l *BatchDeleteCommentsLogic) BatchDeleteComments(in *discussionrpc.BatchDeleteCommentsRequest) (*discussionrpc.BatchDeleteCommentsResponse, error) {
	rows, err := l.svcCtx.TCommentModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchDeleteCommentsResponse{SuccessCount: rows}, nil
}
