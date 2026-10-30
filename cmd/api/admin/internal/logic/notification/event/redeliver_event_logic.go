package event

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type RedeliverEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重新扇出事件
func NewRedeliverEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverEventLogic {
	return &RedeliverEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RedeliverEventLogic) RedeliverEvent(req *types.RedeliverEventReq) (resp *types.EmptyResp, err error) {
	if _, err := l.svcCtx.NotificationService.RedeliverEvent(l.ctx,
		&notificationservice.RedeliverEventRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
