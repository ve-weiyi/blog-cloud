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

type RedeliverEventLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRedeliverEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverEventLogic {
	return &RedeliverEventLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重新扇出事件
func (l *RedeliverEventLogic) RedeliverEvent(in *notificationrpc.RedeliverEventRequest) (*notificationrpc.RedeliverEventResponse, error) {
	err := l.svcCtx.Notify.RedeliverEvent(l.ctx, in.Id)
	if errors.Is(err, notification.ErrEventNotRedeliverable) {
		return nil, bizerr.NewBizError(bizcode.CodeInvalidParam, "事件已完成或已终止，不能重放")
	}
	if err != nil {
		return nil, err
	}
	return &notificationrpc.RedeliverEventResponse{Success: true}, nil
}
