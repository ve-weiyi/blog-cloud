package mqlogic

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
	eventx "github.com/ve-weiyi/vkit/adapter/mqx/event"
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

	bizID := eventx.EventIDFrom(ctx)
	if bizID == "" {
		return mqx.Permanent(errors.New("mq: event id is missing"))
	}

	delivery, err := findDelivery(ctx, l.svcCtx.TNotifyRecordModel, bizID, enums.NotifyChannelSms)
	if err != nil {
		return err
	}
	if delivery != nil && delivery.Status == enums.NotifyRecordStatusSent {
		logger.Infof("短信已发送，跳过重复投递: bizId=%s", bizID)
		return nil
	}

	// 1. 从 SmsProvider 获取真实的模板代码
	templateCode := l.svcCtx.SmsProvider.GetTemplateCode(event.Scene)
	content := ""

	// 2. 尝试从 t_notify_template 表查询模板
	tmpl, err := l.svcCtx.TNotifyTemplateModel.FindOne(ctx, "", "scene = ? AND channel = ? AND enabled = ?", event.Scene, enums.NotifyChannelSms, 1)
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

	// 3. 创建或复用发送记录
	if delivery == nil {
		delivery = &model.TNotifyRecord{
			Channel:      enums.NotifyChannelSms,
			Recipient:    event.Mobile,
			TemplateCode: templateCode,
			Content:      stringToPtr(content),
			Status:       enums.NotifyRecordStatusPending,
			BizId:        bizID,
			CreatedAt:    time.Now(),
		}

		delivery, err = createDelivery(ctx, l.svcCtx.TNotifyRecordModel, delivery)
		if err != nil {
			logger.Errorf("插入短信记录失败: %v", err)
			return err
		}
		if delivery.Status == enums.NotifyRecordStatusSent {
			logger.Infof("短信已发送，跳过重复投递: bizId=%s", bizID)
			return nil
		}
	} else {
		fields := map[string]interface{}{
			"recipient":     event.Mobile,
			"template_code": templateCode,
			"content":       content,
			"status":        enums.NotifyRecordStatusPending,
			"error_msg":     "",
		}
		if _, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id); err != nil {
			logger.Errorf("更新短信记录失败: %v", err)
			return err
		}
	}

	// 4. 发送短信
	sendErr := l.svcCtx.SmsProvider.SendTemplate(ctx, event.Mobile, templateCode, event.Params)

	// 5. 回写发送状态
	fields := make(map[string]interface{})
	if sendErr != nil {
		logger.Errorf("发送短信失败: %v", sendErr)
		fields["status"] = enums.NotifyRecordStatusFailed
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableSmsError(sendErr) {
			logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		logger.Infof("发送短信成功: mobile=%s", event.Mobile)
		fields["status"] = enums.NotifyRecordStatusSent
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
