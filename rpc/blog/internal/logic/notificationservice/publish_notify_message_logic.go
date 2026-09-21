package notificationservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq/mqlogic"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type PublishNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishNotifyMessageLogic {
	return &PublishNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 发布通知消息
func (l *PublishNotifyMessageLogic) PublishNotifyMessage(in *notificationrpc.PublishNotifyMessageRequest) (*notificationrpc.PublishNotifyMessageResponse, error) {
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	if msg.Status != "draft" {
		return &notificationrpc.PublishNotifyMessageResponse{Success: false}, nil
	}

	fields := map[string]interface{}{
		"status":       "published",
		"published_at": time.Now(),
		"updated_at":   time.Now(),
	}

	// WHERE 带上原状态：先 FindById 判状态、再按 id 更新不是原子的。
	// 两次并发发布会同时通过上面那道判断，各自入队一条投递事件。
	affected, err := l.svcCtx.TNotifyMessageModel.UpdateFields(l.ctx, fields, "id = ? AND status = ?", in.Id, "draft")
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		// 已被并发发布，本次不产生副作用
		return &notificationrpc.PublishNotifyMessageResponse{Success: false}, nil
	}

	// 通过 MQ 异步投递 delivery 记录，避免大量用户时接口阻塞
	err = mq.PublishInboxMessageEvent(l.ctx, &mq.InboxMessageEvent{MessageId: msg.Id})
	if errors.Is(err, mq.ErrUnavailable) {
		// MQ 不可用时降级为同步投递。这里复用消费者逻辑，而不是再写一遍扇出：
		// 两份实现必然分叉——此处就曾缺过消费者后来补上的"跳过非已发布消息"判断。
		l.Logger.Infof("消息队列未就绪，降级为同步投递")
		if cerr := mqlogic.NewConsumeInboxMessageLogic(l.svcCtx).Consume(l.ctx, &mq.InboxMessageEvent{MessageId: msg.Id}); cerr != nil {
			l.Logger.Errorf("降级同步投递失败: %v", cerr)
			return nil, cerr
		}
	} else if err != nil {
		l.Logger.Errorf("发送站内信投递消息失败: %v", err)
		return nil, err
	}

	return &notificationrpc.PublishNotifyMessageResponse{
		Success: true,
	}, nil
}
