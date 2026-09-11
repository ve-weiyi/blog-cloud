package mqlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/mail"
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// EmailConsumerLogic Email 消费者逻辑
// 监听邮件消息，发送邮件并记录发送状态
type EmailConsumerLogic struct {
	svcCtx   *svc.ServiceContext
	consumer mqx.Consumer
	logx.Logger
}

// NewEmailConsumerLogic 创建 Email 消费者逻辑
func NewEmailConsumerLogic(svcCtx *svc.ServiceContext) (*EmailConsumerLogic, error) {
	// 类型断言获取 MessageQueue
	emailMQ, ok := mq.EmailMQ.(mqx.MessageQueue)
	if !ok || emailMQ == nil {
		return nil, fmt.Errorf("EmailMQ 未初始化或类型不匹配")
	}

	// 创建消费者
	consumer, err := emailMQ.Consumer(&mqx.ConsumerConfig{
		Topics:        []string{mq.EmailQueue},
		ConsumerName:  "email-consumer",
		PrefetchCount: 10,
		MaxRetries:    3,
	})
	if err != nil {
		return nil, fmt.Errorf("创建邮件消费者失败: %w", err)
	}

	return &EmailConsumerLogic{
		svcCtx:   svcCtx,
		consumer: consumer,
		Logger:   logx.WithContext(context.Background()),
	}, nil
}

// Start 启动消费者
func (l *EmailConsumerLogic) Start() {
	// 订阅主题
	err := l.consumer.Subscribe(mq.EmailRoutingKey)
	if err != nil {
		l.Logger.Errorf("订阅邮件主题失败: %v", err)
		return
	}

	l.Logger.Info("邮件消费者启动成功")

	// 使用处理器消费消息
	err = l.consumer.ConsumeWithHandler(mqx.MessageHandlerFunc(func(ctx context.Context, msg *mqx.Message) error {
		return l.handleMessage(msg.Body)
	}))

	if err != nil {
		l.Logger.Errorf("邮件消费者运行失败: %v", err)
	}
}

// Stop 停止消费者
func (l *EmailConsumerLogic) Stop() {
	ctx := context.Background()
	if err := l.consumer.Stop(ctx); err != nil {
		l.Logger.Errorf("停止邮件消费者失败: %v", err)
	}
}

// handleMessage 处理邮件消息
func (l *EmailConsumerLogic) handleMessage(body []byte) error {
	// 1. 解析消息
	var event mq.EmailMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	l.Logger.Infof("收到邮件消息: email=%s, scene=%s", event.Email, event.Scene)

	// 2. 创建数据库记录
	ctx := context.Background()
	content := ""
	templateCode := ""

	// 2.1 尝试从 t_notify_template 表查询模板
	tmpl, err := l.svcCtx.TNotifyTemplateModel.FindOne(ctx, "", "scene = ? AND channel = ? AND enabled = ?", event.Scene, "email", 1)
	if err == nil {
		if tmpl.Title != "" {
			event.Title = tmpl.Title
		}
		if tmpl.Content != "" {
			content = tmpl.Content
			event.Content = tmpl.Content
		}
		templateCode = tmpl.Code
		l.Logger.Infof("使用邮件模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	} else {
		l.Logger.Infof("未找到邮件模板(scene=%s)，使用默认内容", event.Scene)
		event.Title = "验证码"
		content = fmt.Sprintf("您的验证码是：%s，%s分钟内有效", event.Params["code"], event.Params["time"])
		event.Content = content
	}

	delivery := &model.TNotifyRecord{
		Channel:      "email",
		Recipient:    event.Email,
		TemplateCode: templateCode,
		Content:      stringToPtr(content),
		Status:       "pending",
		BizId:        event.BizId,
		CreatedAt:    time.Now(),
	}

	_, err = l.svcCtx.TNotifyRecordModel.Insert(ctx, delivery)
	if err != nil {
		l.Logger.Errorf("插入邮件记录失败: %v", err)
		return err
	}

	sendErr := l.sendEmail(ctx, &event)

	fields := make(map[string]interface{})
	if sendErr != nil {
		l.Logger.Errorf("发送邮件失败: %v", sendErr)
		fields["status"] = "failed"
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableEmailError(sendErr) {
			l.Logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		l.Logger.Infof("发送邮件成功: email=%s", event.Email)
		fields["status"] = "sent"
		fields["sent_at"] = time.Now()
	}

	_, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
	if err != nil {
		l.Logger.Errorf("更新邮件状态失败: %v", err)
	}

	return sendErr
}

func (l *EmailConsumerLogic) sendEmail(ctx context.Context, event *mq.EmailMessageEvent) error {
	return l.svcCtx.EmailDeliver.DeliveryEmail(&mail.EmailMessage{
		To:      []string{event.Email},
		CC:      []string{},
		Subject: event.Title,
		Content: event.Content,
	})
}

func isNonRetryableEmailError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	nonRetryableErrors := []string{
		"Account is abnormal",
		"service is not open",
		"password is incorrect",
		"invalid email",
		"INVALID_EMAIL",
		"mailbox unavailable",
		"550",
		"553",
	}

	for _, pattern := range nonRetryableErrors {
		if containsStr(errMsg, pattern) {
			return true
		}
	}

	return false
}

func containsStr(s, substr string) bool {
	if s == "" || substr == "" {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
