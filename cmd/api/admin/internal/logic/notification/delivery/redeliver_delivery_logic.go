package delivery

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type RedeliverDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 按投递记录重发
func NewRedeliverDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverDeliveryLogic {
	return &RedeliverDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedeliverDeliveryLogic) RedeliverDelivery(req *types.RedeliverDeliveryReq) (resp *types.EmptyResp, err error) {
	if _, err := l.svcCtx.NotificationService.RedeliverDelivery(l.ctx,
		&notificationservice.RedeliverDeliveryRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
