package mqlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// InboxConsumerLogic 站内信投递消费者逻辑
// 监听消息发布事件，为每个目标用户创建 delivery 记录
type InboxConsumerLogic struct {
	svcCtx   *svc.ServiceContext
	consumer mqx.Consumer
	logx.Logger
}

// NewInboxConsumerLogic 创建站内信投递消费者逻辑
func NewInboxConsumerLogic(svcCtx *svc.ServiceContext) (*InboxConsumerLogic, error) {
	inboxMQ, ok := mq.InboxMQ.(mqx.MessageQueue)
	if !ok || inboxMQ == nil {
		return nil, fmt.Errorf("InboxMQ 未初始化或类型不匹配")
	}

	consumer, err := inboxMQ.Consumer(&mqx.ConsumerConfig{
		Topics:        []string{mq.InboxQueue},
		ConsumerName:  "inbox-consumer",
		PrefetchCount: 10,
		MaxRetries:    3,
	})
	if err != nil {
		return nil, fmt.Errorf("创建站内信消费者失败: %w", err)
	}

	return &InboxConsumerLogic{
		svcCtx:   svcCtx,
		consumer: consumer,
		Logger:   logx.WithContext(context.Background()),
	}, nil
}

// Start 启动消费者
func (l *InboxConsumerLogic) Start() {
	err := l.consumer.Subscribe(mq.InboxRoutingKey)
	if err != nil {
		l.Logger.Errorf("订阅站内信主题失败: %v", err)
		return
	}

	l.Logger.Info("站内信消费者启动成功")

	err = l.consumer.ConsumeWithHandler(mqx.MessageHandlerFunc(func(ctx context.Context, msg *mqx.Message) error {
		return l.handleMessage(msg.Body)
	}))

	if err != nil {
		l.Logger.Errorf("站内信消费者运行失败: %v", err)
	}
}

// Stop 停止消费者
func (l *InboxConsumerLogic) Stop() {
	ctx := context.Background()
	if err := l.consumer.Stop(ctx); err != nil {
		l.Logger.Errorf("停止站内信消费者失败: %v", err)
	}
}

// handleMessage 处理站内信投递消息
func (l *InboxConsumerLogic) handleMessage(body []byte) error {
	// 1. 解析消息
	var event mq.InboxMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	l.Logger.Infof("收到站内信投递消息: messageId=%d", event.MessageId)

	// 2. 查询消息内容
	ctx := context.Background()
	msg, err := l.svcCtx.TNotifyMessageModel.FindById(ctx, event.MessageId)
	if err != nil {
		l.Logger.Errorf("查询消息失败: %v", err)
		return fmt.Errorf("查询消息失败: %v", err)
	}

	// 3. 解析目标用户
	recipients := l.resolveRecipients(msg)
	if len(recipients) == 0 {
		l.Logger.Infof("消息 %d 无目标用户，跳过", event.MessageId)
		return nil
	}

	// 4. 查询已存在的投递记录（幂等性：跳过已投递的用户）
	existing, err := l.svcCtx.TNotifyRecordModel.FindALL(ctx,
		"message_id = ? AND channel = ?", event.MessageId, "inbox")
	if err != nil {
		l.Logger.Errorf("查询已有投递记录失败: %v", err)
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
		l.Logger.Infof("消息 %d 所有用户已投递，跳过", event.MessageId)
		return nil
	}

	_, err = l.svcCtx.TNotifyRecordModel.InsertBatch(ctx, newMessages...)
	if err != nil {
		l.Logger.Errorf("批量插入投递记录失败: %v", err)
		return err
	}

	l.Logger.Infof("消息 %d 投递完成: 目标=%d, 已存在=%d, 新投递=%d",
		event.MessageId, len(recipients), len(existingSet), len(newMessages))
	return nil
}

// resolveRecipients 解析消息的目标用户列表
func (l *InboxConsumerLogic) resolveRecipients(msg *model.TNotifyMessage) []string {
	switch msg.TargetType {
	case "all":
		users, err := l.svcCtx.TUserModel.FindALL(context.Background(), "status = ?", 1)
		if err != nil {
			l.Logger.Errorf("查询全量用户失败: %v", err)
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
