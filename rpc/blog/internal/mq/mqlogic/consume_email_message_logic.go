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
	"github.com/ve-weiyi/vkit/adapter/mail"
)

// ConsumeEmailMessageLogic 邮件消息消费者逻辑
// 监听邮件消息，发送邮件并记录发送状态
type ConsumeEmailMessageLogic struct {
	svcCtx *svc.ServiceContext
}

// NewConsumeEmailMessageLogic 创建邮件消息消费者逻辑
func NewConsumeEmailMessageLogic(svcCtx *svc.ServiceContext) *ConsumeEmailMessageLogic {
	return &ConsumeEmailMessageLogic{svcCtx: svcCtx}
}

// Consume 处理邮件消息
func (l *ConsumeEmailMessageLogic) Consume(ctx context.Context, event *mq.EmailMessageEvent) error {
	logger := logx.WithContext(ctx)
	logger.Infof("收到邮件消息: email=%s, scene=%s", event.Email, event.Scene)

	content := ""
	templateCode := ""

	// 1. 尝试从 t_notify_template 表查询模板
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
		logger.Infof("使用邮件模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	} else {
		logger.Infof("未找到邮件模板(scene=%s)，使用默认内容", event.Scene)
		event.Title = "验证码"
		content = fmt.Sprintf("您的验证码是：%s，%s分钟内有效", event.Params["code"], event.Params["time"])
		event.Content = content
	}

	// 2. 创建发送记录
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
		logger.Errorf("插入邮件记录失败: %v", err)
		return err
	}

	// 3. 发送邮件
	sendErr := l.sendEmail(event)

	// 4. 回写发送状态
	fields := make(map[string]interface{})
	if sendErr != nil {
		logger.Errorf("发送邮件失败: %v", sendErr)
		fields["status"] = "failed"
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableEmailError(sendErr) {
			logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		logger.Infof("发送邮件成功: email=%s", event.Email)
		fields["status"] = "sent"
		fields["sent_at"] = time.Now()
	}

	_, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
	if err != nil {
		logger.Errorf("更新邮件状态失败: %v", err)
	}

	return sendErr
}

// sendEmail 投递邮件。
func (l *ConsumeEmailMessageLogic) sendEmail(event *mq.EmailMessageEvent) error {
	return l.svcCtx.EmailDeliver.DeliveryEmail(&mail.EmailMessage{
		To:      []string{event.Email},
		CC:      []string{},
		Subject: event.Title,
		Content: event.Content,
	})
}

// isNonRetryableEmailError 判断是否为不可重试的错误
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
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
