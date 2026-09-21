package notificationservicelogic

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
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
// 此处**不**在 MQ 不可用时降级为同步投递——本接口是人工补救动作，让调用方稍后重试，
// 比在 RPC 链路里展开全量收件人更安全。
func (l *RedeliverNotifyMessageLogic) RedeliverNotifyMessage(in *notificationrpc.RedeliverNotifyMessageRequest) (*notificationrpc.RedeliverNotifyMessageResponse, error) {
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	// 草稿与已撤回的消息不需要投递，重新投递对它们没有意义
	if msg.Status != "published" {
		return &notificationrpc.RedeliverNotifyMessageResponse{Success: false}, nil
	}

	err = mq.PublishInboxMessageEvent(l.ctx, &mq.InboxMessageEvent{MessageId: msg.Id})
	if errors.Is(err, mq.ErrUnavailable) {
		return nil, errors.New("message queue unavailable, retry later")
	}
	if err != nil {
		return nil, err
	}

	return &notificationrpc.RedeliverNotifyMessageResponse{Success: true}, nil
}
