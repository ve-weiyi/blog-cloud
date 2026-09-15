package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteTalksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteTalksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteTalksLogic {
	return &BatchDeleteTalksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *BatchDeleteTalksLogic) BatchDeleteTalks(in *discussionrpc.BatchDeleteTalksRequest) (*discussionrpc.BatchDeleteTalksResponse, error) {
	rows, err := l.svcCtx.TTalkModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}
	return &discussionrpc.BatchDeleteTalksResponse{SuccessCount: rows}, nil
}
