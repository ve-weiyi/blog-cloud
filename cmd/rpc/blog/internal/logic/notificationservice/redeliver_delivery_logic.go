package notificationservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type RedeliverDeliveryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRedeliverDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverDeliveryLogic {
	return &RedeliverDeliveryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 按投递记录重发
func (l *RedeliverDeliveryLogic) RedeliverDelivery(in *notificationrpc.RedeliverDeliveryRequest) (*notificationrpc.RedeliverDeliveryResponse, error) {
	err := l.svcCtx.Notify.RedeliverDelivery(l.ctx, in.Id)
	if errors.Is(err, notification.ErrDeliveryNotRedeliverable) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "该投递已发送或已取消，不能重发")
	}
	if err != nil {
		return nil, err
	}
	return &notificationrpc.RedeliverDeliveryResponse{Success: true}, nil
}
