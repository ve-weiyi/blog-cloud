package notificationservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
)

var errAlreadyPublished = errors.New("notification already published")

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

	if msg.Status != enums.NoticeStatusDraft {
		return &notificationrpc.PublishNotifyMessageResponse{Success: false}, nil
	}

	now := time.Now()
	fields := map[string]interface{}{
		"status":       enums.NoticeStatusPublished,
		"published_at": now,
		"updated_at":   now,
	}

	// 事件 ID 与信封在这里生成；随后与状态变更同事务写入 outbox。
	event, err := mq.InboxMessage.Message(&mq.InboxMessageEvent{MessageId: msg.Id})
	if err != nil {
		return nil, err
	}

	// WHERE 带上原状态：先 FindById 判状态、再按 id 更新不是原子的。
	// 状态变更与 outbox 事件必须同事务提交，避免「库内已发布、消息从未发出」。
	err = l.svcCtx.GormDB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		affected, err := l.svcCtx.TNotifyMessageModel.WithTx(tx).
			UpdateFields(l.ctx, fields, "id = ? AND status = ?", in.Id, enums.NoticeStatusDraft)
		if err != nil {
			return err
		}
		if affected == 0 {
			return errAlreadyPublished
		}
		return l.svcCtx.Outbox.Enqueue(l.ctx, tx, event)
	})
	if errors.Is(err, errAlreadyPublished) {
		// 已被并发发布，本次不产生副作用
		return &notificationrpc.PublishNotifyMessageResponse{Success: false}, nil
	}
	if err != nil {
		return nil, err
	}

	return &notificationrpc.PublishNotifyMessageResponse{
		Success: true,
	}, nil
}
