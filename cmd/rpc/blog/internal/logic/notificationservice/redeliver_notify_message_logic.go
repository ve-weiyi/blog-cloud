package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

type RedeliverNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRedeliverNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RedeliverNotifyMessageLogic {
	return &RedeliverNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RedeliverNotifyMessage 对指定消息重新投递，用于修补丢失了投递事件的消息。
//
// 幂等：消费者按已有投递记录去重，重复投递不会产生第二条记录，因此可以安全重复调用。
// 事件写入 outbox 后由 relay 负责投递，接口本身不依赖 broker 实时可用。
func (l *RedeliverNotifyMessageLogic) RedeliverNotifyMessage(in *notificationrpc.RedeliverNotifyMessageRequest) (*notificationrpc.RedeliverNotifyMessageResponse, error) {
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 草稿与已撤回的消息不需要投递，重新投递对它们没有意义
	if msg.Status != enums.NoticeStatusPublished {
		return &notificationrpc.RedeliverNotifyMessageResponse{Success: false}, nil
	}

	event, err := mq.InboxMessage.Message(&mq.InboxMessageEvent{MessageId: msg.Id})
	if err != nil {
		return nil, err
	}
	// 重投没有其他写入需要原子化，落库本身就是本次动作的全部副作用
	if err := l.svcCtx.Outbox.Enqueue(l.ctx, l.svcCtx.GormDB, event); err != nil {
		return nil, err
	}

	return &notificationrpc.RedeliverNotifyMessageResponse{Success: true}, nil
}
