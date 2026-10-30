package consumer

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

// ConsumeNotificationFanoutLogic 通知扇出消费者：把事件展开成收件箱条目。
//
// 消息只带事件 id，载荷与游标都在库里，所以重复投递是安全的：从游标继续即可。
type ConsumeNotificationFanoutLogic struct {
	svcCtx *svc.ServiceContext
}

func NewConsumeNotificationFanoutLogic(svcCtx *svc.ServiceContext) *ConsumeNotificationFanoutLogic {
	return &ConsumeNotificationFanoutLogic{svcCtx: svcCtx}
}

// Consume 处理一条扇出唤醒消息。
func (l *ConsumeNotificationFanoutLogic) Consume(ctx context.Context, event *mq.NotificationFanoutEvent) error {
	return l.svcCtx.Notify.Fanout(ctx, event.EventID)
}
