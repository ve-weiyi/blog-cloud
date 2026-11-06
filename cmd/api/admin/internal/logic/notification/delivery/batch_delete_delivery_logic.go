package delivery

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type BatchDeleteDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除投递记录
func NewBatchDeleteDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteDeliveryLogic {
	return &BatchDeleteDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteDeliveryLogic) BatchDeleteDelivery(req *types.IdsReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.BatchDeleteDeliveries(l.ctx,
		&notificationservice.BatchDeleteDeliveriesRequest{Ids: req.Ids})
	if err != nil {
		return nil, err
	}
	return &types.BatchResp{SuccessCount: out.SuccessCount}, nil
}
