package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchUpdateCommentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchUpdateCommentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateCommentStatusLogic {
	return &BatchUpdateCommentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量更新评论状态
func (l *BatchUpdateCommentStatusLogic) BatchUpdateCommentStatus(in *discussionrpc.BatchUpdateCommentStatusRequest) (*discussionrpc.BatchUpdateCommentStatusResponse, error) {
	rows, err := l.svcCtx.TCommentModel.UpdateFields(l.ctx, map[string]interface{}{
		"status": in.Status,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchUpdateCommentStatusResponse{SuccessCount: rows}, nil
}
