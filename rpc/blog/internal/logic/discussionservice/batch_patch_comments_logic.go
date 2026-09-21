package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchPatchCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchPatchCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchPatchCommentsLogic {
	return &BatchPatchCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量更新评论状态
func (l *BatchPatchCommentsLogic) BatchPatchComments(in *discussionrpc.BatchPatchCommentsRequest) (*discussionrpc.BatchPatchCommentsResponse, error) {
	rows, err := l.svcCtx.TCommentModel.UpdateFields(l.ctx, map[string]interface{}{
		"status": in.Status,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchPatchCommentsResponse{SuccessCount: rows}, nil
}
