package mqlogic

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/mailtemplate"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/adapter/emailx"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
	eventx "github.com/ve-weiyi/vkit/adapter/mqx/event"
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

	bizID := eventx.EventIDFrom(ctx)
	if bizID == "" {
		return mqx.Permanent(errors.New("mq: event id is missing"))
	}

	delivery, err := findDelivery(ctx, l.svcCtx.TNotifyRecordModel, bizID, enums.NotifyChannelEmail)
	if err != nil {
		return err
	}
	if delivery != nil && delivery.Status == enums.NotifyRecordStatusSent {
		logger.Infof("邮件已发送，跳过重复投递: bizId=%s", bizID)
		return nil
	}

	content := ""
	templateCode := ""

	// 1. 尝试从 t_notify_template 表查询模板
	tmpl, err := l.svcCtx.TNotifyTemplateModel.FindOne(ctx, "", "scene = ? AND channel = ? AND enabled = ?", event.Scene, enums.NotifyChannelEmail, 1)
	if err == nil {
		if tmpl.Title != "" {
			event.Title = tmpl.Title
		}
		if tmpl.Content != "" {
			title, intro := mailtemplate.CodeSceneText(event.Scene)
			if tmpl.Title != "" {
				title = tmpl.Title
			}
			params := make(map[string]string, len(event.Params)+2)
			for key, value := range event.Params {
				params[key] = value
			}
			params["title"] = title
			params["intro"] = intro

			rendered, renderErr := mailtemplate.Render(tmpl.Content, params)
			if renderErr != nil {
				logger.Errorf("渲染邮件模板失败: %v", renderErr)
				rendered, renderErr = mailtemplate.CodeEmail(event.Scene, event.Params["code"], event.Params["time"])
				if renderErr != nil {
					return renderErr
				}
				event.Title = title
			}
			content = rendered
			event.Content = content
		}
		templateCode = tmpl.Code
		logger.Infof("使用邮件模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	} else {
		logger.Infof("未找到邮件模板(scene=%s)，使用默认内容", event.Scene)
		rendered, renderErr := mailtemplate.CodeEmail(event.Scene, event.Params["code"], event.Params["time"])
		if renderErr != nil {
			return renderErr
		}
		event.Title, _ = mailtemplate.CodeSceneText(event.Scene)
		content = rendered
		event.Content = content
	}

	// 2. 创建或复用发送记录
	if delivery == nil {
		delivery = &model.TNotifyRecord{
			Channel:      enums.NotifyChannelEmail,
			Recipient:    event.Email,
			TemplateCode: templateCode,
			Content:      stringToPtr(content),
			Status:       enums.NotifyRecordStatusPending,
			BizId:        bizID,
			CreatedAt:    time.Now(),
		}

		delivery, err = createDelivery(ctx, l.svcCtx.TNotifyRecordModel, delivery)
		if err != nil {
			logger.Errorf("插入邮件记录失败: %v", err)
			return err
		}
		if delivery.Status == enums.NotifyRecordStatusSent {
			logger.Infof("邮件已发送，跳过重复投递: bizId=%s", bizID)
			return nil
		}
	} else {
		fields := map[string]interface{}{
			"recipient":     event.Email,
			"template_code": templateCode,
			"content":       content,
			"status":        enums.NotifyRecordStatusPending,
			"error_msg":     "",
		}
		if _, err = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id); err != nil {
			logger.Errorf("更新邮件记录失败: %v", err)
			return err
		}
	}

	// 3. 发送邮件
	sendErr := l.sendEmail(event)

	// 4. 回写发送状态
	fields := make(map[string]interface{})
	if sendErr != nil {
		logger.Errorf("发送邮件失败: %v", sendErr)
		fields["status"] = enums.NotifyRecordStatusFailed
		fields["error_msg"] = sendErr.Error()

		if isNonRetryableEmailError(sendErr) {
			logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = l.svcCtx.TNotifyRecordModel.UpdateFields(ctx, fields, "id = ?", delivery.Id)
			return nil
		}
	} else {
		logger.Infof("发送邮件成功: email=%s", event.Email)
		fields["status"] = enums.NotifyRecordStatusSent
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
	return l.svcCtx.EmailDeliver.DeliveryEmail(&emailx.EmailMessage{
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
