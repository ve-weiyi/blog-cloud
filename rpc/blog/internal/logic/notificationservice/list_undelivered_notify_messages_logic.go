package notificationservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

// 查询参数默认值与上限
const (
	defaultUndeliveredBeforeSeconds = 60
	defaultUndeliveredLimit         = 100
	maxUndeliveredLimit             = 500
)

type ListUndeliveredNotifyMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUndeliveredNotifyMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUndeliveredNotifyMessagesLogic {
	return &ListUndeliveredNotifyMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListUndeliveredNotifyMessages 查询已发布但尚未生成任何投递记录的消息。
//
// 投递记录由消费者在投递时生成，因此这类消息对读者不可见、也拿不到任何报错——
// 它们只可能来自"发布与投递之间中断"（进程在写状态与发投递事件之间退出）。
// 本查询是这种静默失效的唯一现成排查入口。
func (l *ListUndeliveredNotifyMessagesLogic) ListUndeliveredNotifyMessages(in *notificationrpc.ListUndeliveredNotifyMessagesRequest) (*notificationrpc.ListUndeliveredNotifyMessagesResponse, error) {
	beforeSeconds := in.BeforeSeconds
	if beforeSeconds <= 0 {
		beforeSeconds = defaultUndeliveredBeforeSeconds
	}

	limit := in.Limit
	if limit <= 0 {
		limit = defaultUndeliveredLimit
	}
	if limit > maxUndeliveredLimit {
		limit = maxUndeliveredLimit
	}

	// 只看过了宽限期仍未投递的：发布后立即查询会把正在投递中的消息误报出来
	cutoff := time.Now().Add(-time.Duration(beforeSeconds) * time.Second)

	// NOT EXISTS 走 t_notify_record 的 message_id 索引，不需要把记录拉回来比对
	conditions := "status = ? AND published_at IS NOT NULL AND published_at < ?" +
		" AND NOT EXISTS (SELECT 1 FROM `t_notify_record` `r`" +
		" WHERE `r`.`message_id` = `t_notify_message`.`id` AND `r`.`channel` = ?)"

	messages, _, err := l.svcCtx.TNotifyMessageModel.FindListAndTotal(
		l.ctx, 1, int(limit), "published_at asc", conditions, "published", cutoff, "inbox")
	if err != nil {
		return nil, err
	}

	list := make([]*notificationrpc.NotifyMessage, 0, len(messages))
	for _, m := range messages {
		list = append(list, convertTNotifyMessageToProto(m))
	}

	return &notificationrpc.ListUndeliveredNotifyMessagesResponse{
		List: list,
	}, nil
}
