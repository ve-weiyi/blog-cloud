package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteMessagesLogic {
	return &BatchDeleteMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除留言
func (l *BatchDeleteMessagesLogic) BatchDeleteMessages(in *discussionrpc.BatchDeleteMessagesRequest) (*discussionrpc.BatchDeleteMessagesResponse, error) {
	rows, err := l.svcCtx.TMessageModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &discussionrpc.BatchDeleteMessagesResponse{SuccessCount: rows}, nil
}
