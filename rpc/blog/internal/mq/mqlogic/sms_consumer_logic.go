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
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// SmsConsumerLogic SMS 消费者逻辑
// 监听短信消息，发送短信并记录发送状态
type SmsConsumerLogic struct {
	svcCtx   *svc.ServiceContext
	consumer mqx.Consumer
	logx.Logger
}

// NewSmsConsumerLogic 创建 SMS 消费者逻辑
func NewSmsConsumerLogic(svcCtx *svc.ServiceContext) (*SmsConsumerLogic, error) {
	// 类型断言获取 MessageQueue
	smsMQ, ok := mq.SmsMQ.(mqx.MessageQueue)
	if !ok || smsMQ == nil {
		return nil, fmt.Errorf("SmsMQ 未初始化或类型不匹配")
	}

	// 创建消费者
	consumer, err := smsMQ.Consumer(&mqx.ConsumerConfig{
		Topics:        []string{mq.SmsQueue},
		ConsumerName:  "sms-consumer",
		PrefetchCount: 10,
		MaxRetries:    3,
	})
	if err != nil {
		return nil, fmt.Errorf("创建短信消费者失败: %w", err)
	}

	return &SmsConsumerLogic{
		svcCtx:   svcCtx,
		consumer: consumer,
		Logger:   logx.WithContext(context.Background()),
	}, nil
}

func (l *SmsConsumerLogic) Start() {
	// 订阅主题
	err := l.consumer.Subscribe(mq.SmsRoutingKey)
	if err != nil {
		l.Logger.Errorf("订阅短信主题失败: %v", err)
		return
	}

	l.Logger.Info("短信消费者启动成功")

	// 使用处理器消费消息
	err = l.consumer.ConsumeWithHandler(mqx.MessageHandlerFunc(func(ctx context.Context, msg *mqx.Message) error {
		return l.handleMessage(msg.Body)
	}))

	if err != nil {
		l.Logger.Errorf("短信消费者运行失败: %v", err)
	}
}

func (l *SmsConsumerLogic) Stop() {
	ctx := context.Background()
	if err := l.consumer.Stop(ctx); err != nil {
		l.Logger.Errorf("停止短信消费者失败: %v", err)
	}
}

func (l *SmsConsumerLogic) handleMessage(body []byte) error {
	// 1. 解析消息
	var event mq.SmsMessageEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	l.Logger.Infof("收到短信消息: mobile=%s, scene=%s", event.Mobile, event.Scene)

	// 2. 从 SmsProvider 获取真实的模板代码
	templateCode := l.svcCtx.SmsProvider.GetTemplateCode(event.Scene)
	ctx := context.Background()
	content := ""

	// 2.1 尝试从 t_notify_template 表查询模板
	tmpl, err := l.svcCtx.TNotifyTemplateModel.FindOne(ctx, "", "scene = ? AND channel = ? AND enabled = ?", event.Scene, "sms", 1)
	if err == nil {
		if tmpl.Content != "" {
			content = tmpl.Content
		}
		if tmpl.Code != "" {
			templateCode = tmpl.Code
		}
		l.Logger.Infof("使用短信模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	} else {
		l.Logger.Infof("未找到短信模板(scene=%s)，使用默认内容", event.Scene)
		content = fmt.Sprintf("您的验证码是：%s，%s分钟内有效", event.Params["code"], event.Params["time"])
	}

	// 3. 创建数据库记录
	delivery := &model.TNotifyRecord{
		Channel:      "sms",
		Recipient:    event.Mobile,
		TemplateCode: templateCode,
		Content:      stringToPtr(content),
		Status:       "pending",
		BizId:        event.BizId,
		CreatedAt:    time.Now(),
	}

	// 3. 插入数据库
	_, err = l.svcCtx.TNotifyRecordModel.Insert(ctx, delivery)
	if err != nil {
		l.Logger.Errorf("插入短信记录失败: %v", err)
		return err
	}

	// 4. 发送短信
	sendErr := l.sendSms(ctx, &event, templateCode)

	// 5. 更新发送状态
	fields := make(map[string]interface{})
	if sendErr != nil {
		l.Logger.Errorf("发送短信失败: %v", sendErr)
		fields["status"] = "failed"
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableError(sendErr) {
			l.Logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		l.Logger.Infof("发送短信成功: mobile=%s", event.Mobile)
		fields["status"] = "sent"
		fields["sent_at"] = time.Now()
	}

	_, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
	if err != nil {
		l.Logger.Errorf("更新短信状态失败: %v", err)
	}

	return sendErr
}

func (l *SmsConsumerLogic) sendSms(ctx context.Context, event *mq.SmsMessageEvent, templateCode string) error {

	// 否则使用模板发送
	return l.svcCtx.SmsProvider.SendTemplate(ctx, event.Mobile, templateCode, event.Params)
}

// isNonRetryableError 判断是否为不可重试的错误
func isNonRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	nonRetryableErrors := []string{
		"MOBILE_NUMBER_ILLEGAL",
		"INVALID_PHONE_NUMBER",
		"Account is abnormal",
		"service is not open",
		"password is incorrect",
		"INVALID_PARAMETERS",
		"TEMPLATE_MISSING_PARAMETERS",
		"SMS_TEMPLATE_ILLEGAL",
		"isv.SMS_TEMPLATE_ILLEGAL",
	}

	for _, pattern := range nonRetryableErrors {
		if contains(errMsg, pattern) {
			return true
		}
	}

	return false
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	if s == "" || substr == "" {
		return false
	}
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
