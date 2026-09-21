package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type RedeliverNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重新投递通知消息
func NewRedeliverNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverNotifyMessageLogic {
	return &RedeliverNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RedeliverNotifyMessage 对指定消息重新投递。
//
// 幂等：消费者按已有投递记录去重，重复调用不会产生第二条记录。
// 消息不是已发布状态时返回 0，表示本次未产生副作用。
func (l *RedeliverNotifyMessageLogic) RedeliverNotifyMessage(req *types.RedeliverNotifyMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.RedeliverNotifyMessage(l.ctx, &notificationservice.RedeliverNotifyMessageRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	var successCount int64
	if out.Success {
		successCount = 1
	}

	return &types.BatchResp{
		SuccessCount: successCount,
	}, nil
}
