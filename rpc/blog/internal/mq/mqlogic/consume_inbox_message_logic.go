package mqlogic

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"

	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
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

	// 2. 只有已发布的消息才展开投递。
	// 事件从发出到被消费之间消息状态可能已经变化（撤回是最常见的一种），而撤回路径只回收
	// "当时已存在"的未读记录——若在这里不拦，为一条已撤回的消息新插入的记录不会被回收，
	// 公告会照常出现在所有人的未读列表里。
	// 这里返回 nil 而不是错误：这是正常的业务终态，不该触发消息重投。
	if msg.Status != "published" {
		logger.Infof("消息 %d 当前状态为 %s，跳过投递", event.MessageId, msg.Status)
		return nil
	}

	// 3. 解析目标用户
	recipients := l.resolveRecipients(ctx, msg)
	if len(recipients) == 0 {
		logger.Infof("消息 %d 无目标用户，跳过", event.MessageId)
		return nil
	}

	// 4. 查询已存在的投递记录（幂等性：跳过已投递的用户）
	existing, err := l.svcCtx.TNotifyRecordModel.FindALL(ctx,
		"message_id = ? AND channel = ?", event.MessageId, "inbox")
	if err != nil {
		logger.Errorf("查询已有投递记录失败: %v", err)
	}

	existingSet := make(map[string]bool, len(existing))
	for _, d := range existing {
		existingSet[d.Recipient] = true
	}

	// 5. 过滤出新用户，批量插入
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

	// 投递记录已落库，此刻广播事件才安全：客户端收到信号后会立即回查未读，
	// 若信号早于记录写入，那次回查就查不到，通知要等下一次信号才可见
	publishedAt := now.Unix()
	if msg.PublishedAt != nil {
		publishedAt = msg.PublishedAt.Unix()
	}
	if err := notifyx.Publish(ctx, l.svcCtx.Redis, notifyx.Event{
		Type:        notifyx.EventNotice,
		MessageId:   msg.Id,
		Title:       msg.Title,
		Category:    msg.Category,
		Level:       msg.Level,
		PublishedAt: publishedAt,
	}); err != nil {
		logger.Errorf("广播通知发布事件失败: %v", err)
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
