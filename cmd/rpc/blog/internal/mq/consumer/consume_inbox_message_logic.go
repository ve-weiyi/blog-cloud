package consumer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/model"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/vkit/x/jsonv"

	"github.com/ve-weiyi/blog-cloud/infra/messaging/notifyx"
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

	var (
		msg             model.TNotifyMessage
		recipientCount  int
		existingCount   int
		newCount        int
		shouldBroadcast bool
	)

	// 消息行锁把「读取状态」与「插入投递记录」串行化：撤回事务会等待本事务提交，
	// 之后再回收本次插入的 unread 记录；若撤回先提交，这里读到 revoked 直接跳过。
	err := l.svcCtx.GormDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("`t_notify_message`").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("`id` = ?", event.MessageId).
			First(&msg).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}

		// 只有已发布的消息才展开投递。
		// 事件从发出到被消费之间消息状态可能已经变化（撤回是最常见的一种），而撤回路径只回收
		// "当时已存在"的未读记录——若在这里不拦，为一条已撤回的消息新插入的记录不会被回收，
		// 公告会照常出现在所有人的未读列表里。
		// 这里返回 nil 而不是错误：这是正常的业务终态，不该触发消息重投。
		if msg.Status != enums.NoticeStatusPublished {
			logger.Infof("消息 %d 当前状态为 %s，跳过投递", event.MessageId, msg.Status)
			return nil
		}

		recipients, err := l.resolveRecipients(ctx, &msg)
		if err != nil {
			return err
		}
		recipientCount = len(recipients)
		if recipientCount == 0 {
			logger.Infof("消息 %d 无目标用户，跳过", event.MessageId)
			return nil
		}

		// 幂等：跳过已投递的用户。查询与插入同处一个事务，避免查重窗口。
		existing, err := l.svcCtx.TNotifyRecordModel.WithTx(tx).FindAll(ctx,
			"message_id = ? AND channel = ?", event.MessageId, enums.NotifyChannelInbox)
		if err != nil {
			return err
		}
		existingCount = len(existing)

		existingSet := make(map[string]bool, len(existing))
		for _, d := range existing {
			existingSet[d.Recipient] = true
		}

		now := time.Now()
		var newMessages []*model.TNotifyRecord
		for _, userId := range recipients {
			if existingSet[userId] {
				continue
			}
			newMessages = append(newMessages, &model.TNotifyRecord{
				MessageId: msg.Id,
				Channel:   enums.NotifyChannelInbox,
				Recipient: userId,
				Content:   msg.Content,
				Status:    enums.NotifyRecordStatusUnread,
				BizId:     fmt.Sprintf("%d:%s", msg.Id, userId),
				CreatedAt: now,
			})
		}

		if len(newMessages) == 0 {
			logger.Infof("消息 %d 所有用户已投递，跳过", event.MessageId)
			return nil
		}

		if _, err := l.svcCtx.TNotifyRecordModel.WithTx(tx).InsertBatch(ctx, newMessages...); err != nil {
			return err
		}

		newCount = len(newMessages)
		shouldBroadcast = true
		return nil
	})
	if err != nil {
		logger.Errorf("站内信投递失败: messageId=%d err=%v", event.MessageId, err)
		return err
	}
	if !shouldBroadcast {
		return nil
	}

	// 投递记录已提交，此刻广播事件才安全：客户端收到信号后会立即回查未读，
	// 若信号早于记录写入，那次回查就查不到，通知要等下一次信号才可见
	publishedAt := time.Now().Unix()
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
		event.MessageId, recipientCount, existingCount, newCount)
	return nil
}

// resolveRecipients 解析消息的目标用户列表。
// 查询失败必须上抛：当成「无目标用户」会 ack 掉消息，整条公告静默丢失。
func (l *ConsumeInboxMessageLogic) resolveRecipients(ctx context.Context, msg *model.TNotifyMessage) ([]string, error) {
	switch msg.TargetType {
	case enums.NotifyTargetTypeAll:
		users, err := l.svcCtx.TUserModel.FindAll(ctx, "status = ?", enums.UserStatusNormal)
		if err != nil {
			return nil, fmt.Errorf("查询全量用户: %w", err)
		}
		ids := make([]string, 0, len(users))
		for _, u := range users {
			ids = append(ids, u.UserId)
		}
		return ids, nil
	case enums.NotifyTargetTypeUserIDs:
		if msg.TargetIds == nil || *msg.TargetIds == "" {
			return nil, nil
		}
		var raw []string
		if err := jsonv.JsonToAny(*msg.TargetIds, &raw); err != nil {
			return nil, nil
		}
		ids := make([]string, 0, len(raw))
		seen := make(map[string]struct{}, len(raw))
		for _, id := range raw {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
		return ids, nil
	default:
		return nil, nil
	}
}
