package mqlogic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// ConsumeSmsMessageLogic SMS 消费者逻辑
// 监听短信消息，发送短信并记录发送状态
type ConsumeSmsMessageLogic struct {
	svcCtx *svc.ServiceContext
}

// NewConsumeSmsMessageLogic 创建 SMS 消费者逻辑
func NewConsumeSmsMessageLogic(svcCtx *svc.ServiceContext) *ConsumeSmsMessageLogic {
	return &ConsumeSmsMessageLogic{svcCtx: svcCtx}
}

// Consume 处理短信消息
func (l *ConsumeSmsMessageLogic) Consume(ctx context.Context, event *mq.SmsMessageEvent) error {
	logger := logx.WithContext(ctx)
	logger.Infof("收到短信消息: mobile=%s, scene=%s", event.Mobile, event.Scene)

	// 1. 从 SmsProvider 获取真实的模板代码
	templateCode := l.svcCtx.SmsProvider.GetTemplateCode(event.Scene)
	content := ""

	// 2. 尝试从 t_notify_template 表查询模板
	tmpl, err := l.svcCtx.TNotifyTemplateModel.FindOne(ctx, "", "scene = ? AND channel = ? AND enabled = ?", event.Scene, "sms", 1)
	if err == nil {
		if tmpl.Content != "" {
			content = tmpl.Content
		}
		if tmpl.Code != "" {
			templateCode = tmpl.Code
		}
		logger.Infof("使用短信模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	} else {
		logger.Infof("未找到短信模板(scene=%s)，使用默认内容", event.Scene)
		content = fmt.Sprintf("您的验证码是：%s，%s分钟内有效", event.Params["code"], event.Params["time"])
	}

	// 3. 创建发送记录
	delivery := &model.TNotifyRecord{
		Channel:      "sms",
		Recipient:    event.Mobile,
		TemplateCode: templateCode,
		Content:      stringToPtr(content),
		Status:       "pending",
		BizId:        event.BizId,
		CreatedAt:    time.Now(),
	}

	_, err = l.svcCtx.TNotifyRecordModel.Insert(ctx, delivery)
	if err != nil {
		logger.Errorf("插入短信记录失败: %v", err)
		return err
	}

	// 4. 发送短信
	sendErr := l.svcCtx.SmsProvider.SendTemplate(ctx, event.Mobile, templateCode, event.Params)

	// 5. 回写发送状态
	fields := make(map[string]interface{})
	if sendErr != nil {
		logger.Errorf("发送短信失败: %v", sendErr)
		fields["status"] = "failed"
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableSmsError(sendErr) {
			logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		logger.Infof("发送短信成功: mobile=%s", event.Mobile)
		fields["status"] = "sent"
		fields["sent_at"] = time.Now()
	}

	_, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
	if err != nil {
		logger.Errorf("更新短信状态失败: %v", err)
	}

	return sendErr
}

// isNonRetryableSmsError 判断是否为不可重试的错误
func isNonRetryableSmsError(err error) bool {
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
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
