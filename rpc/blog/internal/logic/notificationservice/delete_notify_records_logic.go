package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"

	notificationrpc2 "github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteNotifyRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNotifyRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotifyRecordsLogic {
	return &DeleteNotifyRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteNotifyRecordsLogic) DeleteNotifyRecords(in *notificationrpc2.DeleteNotifyRecordsRequest) (*notificationrpc2.DeleteNotifyRecordsResponse, error) {
	var successCount int64
	for _, id := range in.Ids {
		rows, err := l.svcCtx.TNotifyRecordModel.Delete(l.ctx, id)
		if err != nil {
			return nil, err
		}
		successCount += rows
	}

	return &notificationrpc.DeleteNotifyRecordsResponse{
		SuccessCount: successCount,
	}, nil
}
