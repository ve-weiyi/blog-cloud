package mqlogic

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// ConsumeInboxMessageLogic 站内信投递消费者逻辑
// 监听消息发布事件，为每个目标用户创建 delivery 记录
type ConsumeInboxMessageLogic struct {
	svcCtx *svc.ServiceContext
}

// NewConsumeInboxMessageLogic 创建站内信投递消费者逻辑
func NewConsumeInboxMessageLogic(svcCtx *svc.ServiceContext) *ConsumeInboxMessageLogic {
	return &ConsumeInboxMessageLogic{svcCtx: svcCtx}
}

// Consume 处理站内信投递消息
func (l *ConsumeInboxMessageLogic) Consume(ctx context.Context, event *mq.InboxMessageEvent) error {
	logger := logx.WithContext(ctx)
	logger.Infof("收到站内信投递消息: messageId=%d", event.MessageId)

	// 1. 查询消息内容
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(ctx, event.MessageId)
	if err != nil {
		logger.Errorf("查询消息失败: %v", err)
		return err
	}

	// 2. 解析目标用户
	recipients := l.resolveRecipients(ctx, msg)
	if len(recipients) == 0 {
		logger.Infof("消息 %d 无目标用户，跳过", event.MessageId)
		return nil
	}

	// 3. 查询已存在的投递记录（幂等性：跳过已投递的用户）
	existing, err := l.svcCtx.TNotifyRecordModel.FindALL(ctx,
		"message_id = ? AND channel = ?", event.MessageId, "inbox")
	if err != nil {
		logger.Errorf("查询已有投递记录失败: %v", err)
	}

	existingSet := make(map[string]bool, len(existing))
	for _, d := range existing {
		existingSet[d.Recipient] = true
	}

	// 4. 过滤出新用户，批量插入
	var newMessages []*model.TNotifyRecord
	now := time.Now()
	for _, userId := range recipients {
		if existingSet[userId] {
			continue
		}
		newMessages = append(newMessages, &model.TNotifyRecord{
			MessageId: msg.Id,
			Channel:   "inbox",
			Recipient: userId,
			Content:   msg.Content,
			Status:    "unread",
			CreatedAt: now,
		})
	}

	if len(newMessages) == 0 {
		logger.Infof("消息 %d 所有用户已投递，跳过", event.MessageId)
		return nil
	}

	_, err = l.svcCtx.TNotifyRecordModel.InsertBatch(ctx, newMessages...)
	if err != nil {
		logger.Errorf("批量插入投递记录失败: %v", err)
		return err
	}

	logger.Infof("消息 %d 投递完成: 目标=%d, 已存在=%d, 新投递=%d",
		event.MessageId, len(recipients), len(existingSet), len(newMessages))
	return nil
}

// resolveRecipients 解析消息的目标用户列表
func (l *ConsumeInboxMessageLogic) resolveRecipients(ctx context.Context, msg *model.TNotifyMessage) []string {
	switch msg.TargetType {
	case "all":
		users, err := l.svcCtx.TUserModel.FindALL(ctx, "status = ?", 1)
		if err != nil {
			logx.WithContext(ctx).Errorf("查询全量用户失败: %v", err)
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
