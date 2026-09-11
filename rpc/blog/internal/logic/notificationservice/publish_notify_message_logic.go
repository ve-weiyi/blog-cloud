package notificationservicelogic

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/x/jsonconv"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/notificationrpc"
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

	_, err = l.svcCtx.TNotifyMessageModel.UpdateFields(l.ctx, fields, "id = ?", in.Id)
	if err != nil {
		return nil, err
	}

	// 通过 MQ 异步投递 delivery 记录，避免大量用户时接口阻塞
	if mq.InboxProducer != nil {
		event := &mq.InboxMessageEvent{MessageId: msg.Id}
		err = mq.InboxProducer.Send(l.ctx, &mqx.Message{
			Topic:     mq.InboxQueue,
			Key:       mq.InboxRoutingKey,
			Body:      []byte(jsonconv.AnyToJsonNE(event)),
			Timestamp: time.Now(),
		})
		if err != nil {
			l.Logger.Errorf("发送站内信投递消息失败: %v", err)
			return nil, err
		}
	} else {
		// MQ 不可用时降级为同步投递
		l.Logger.Infof("InboxProducer 未初始化，使用同步投递")
		recipients := l.resolveRecipients(msg)
		now := time.Now()
		for _, userId := range recipients {
			delivery := &model.TNotifyRecord{
				MessageId: msg.Id,
				Channel:   "inbox",
				Recipient: userId,
				Content:   msg.Content,
				Status:    "unread",
				CreatedAt: now,
			}
			_, err = l.svcCtx.TNotifyRecordModel.Insert(l.ctx, delivery)
			if err != nil {
				l.Errorf("PublishNotifyMessage Insert delivery error: %v", err)
			}
		}
	}

	return &notificationrpc.PublishNotifyMessageResponse{
		Success: true,
	}, nil
}

func (l *PublishNotifyMessageLogic) resolveRecipients(msg *model.TNotifyMessage) []string {
	switch msg.TargetType {
	case "all":
		users, err := l.svcCtx.TUserModel.FindALL(l.ctx, "status = ?", 1)
		if err != nil {
			return nil
		}
		ids := make([]string, 0, len(users))
		for _, u := range users {
			ids = append(ids, u.UserId)
		}
		return ids
	case "user_ids":
		if msg.TargetIds == nil || *msg.TargetIds == "" {
			return nil
		}
		return strings.Split(*msg.TargetIds, ",")
	default:
		return nil
	}
}
