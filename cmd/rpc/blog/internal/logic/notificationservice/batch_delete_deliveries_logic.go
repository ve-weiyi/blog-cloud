package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type BatchDeleteDeliveriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteDeliveriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteDeliveriesLogic {
	return &BatchDeleteDeliveriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除投递记录
func (l *BatchDeleteDeliveriesLogic) BatchDeleteDeliveries(in *notificationrpc.BatchDeleteDeliveriesRequest) (*notificationrpc.BatchDeleteDeliveriesResponse, error) {
	if len(in.Ids) == 0 {
		return &notificationrpc.BatchDeleteDeliveriesResponse{}, nil
	}

	rows, err := l.svcCtx.TNotificationDeliveryModel.DeleteBatch(l.ctx, "id IN ?", in.Ids)
	if err != nil {
		return nil, err
	}
	return &notificationrpc.BatchDeleteDeliveriesResponse{SuccessCount: rows}, nil
}
