package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteNotifyRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteNotifyRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyRecordsLogic {
	return &BatchDeleteNotifyRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDeleteNotifyRecordsLogic) BatchDeleteNotifyRecords(in *notificationrpc.BatchDeleteNotifyRecordsRequest) (*notificationrpc.BatchDeleteNotifyRecordsResponse, error) {
	var successCount int64
	for _, id := range in.Ids {
		rows, err := l.svcCtx.TNotifyRecordModel.Delete(l.ctx, id)
		if err != nil {
			return nil, err
		}
		successCount += rows
	}

	return &notificationrpc.BatchDeleteNotifyRecordsResponse{
		SuccessCount: successCount,
	}, nil
}
