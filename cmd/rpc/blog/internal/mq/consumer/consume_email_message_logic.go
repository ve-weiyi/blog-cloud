package consumer

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/mailtemplate"
	"github.com/ve-weiyi/vkit/adapter/emailx"
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
	logx.WithContext(ctx).Infof("收到邮件消息: email=%s, scene=%s", event.Email, event.Scene)

	return (&delivery{
		records: l.svcCtx.TNotifyRecordModel,
		ch:      &emailChannel{svcCtx: l.svcCtx, event: event},
		bizID:   eventx.EventIDFrom(ctx),
	}).run(ctx)
}

// emailChannel 把邮件事件适配到统一投递流程。
type emailChannel struct {
	svcCtx *svc.ServiceContext
	event  *mq.EmailMessageEvent
}

func (c *emailChannel) name() string      { return enums.NotifyChannelEmail }
func (c *emailChannel) recipient() string { return c.event.Email }

// retryable 把厂商文案判定的「不可重试」翻转成「是否值得重试」。
func (c *emailChannel) retryable(err error) bool { return !isNonRetryableEmailError(err) }

// resolve 优先使用库里的模板，缺失时回落到内置文案。
func (c *emailChannel) resolve(ctx context.Context) (string, string, error) {
	logger := logx.WithContext(ctx)

	tmpl, err := c.svcCtx.TNotifyTemplateModel.FindOne(ctx, "",
		"scene = ? AND channel = ? AND enabled = ?", c.event.Scene, enums.NotifyChannelEmail, 1)
	if err != nil {
		logger.Infof("未找到邮件模板(scene=%s)，使用默认内容", c.event.Scene)
		rendered, renderErr := mailtemplate.CodeEmail(c.event.Scene, c.event.Params["code"], c.event.Params["time"])
		if renderErr != nil {
			return "", "", renderErr
		}
		c.event.Title, _ = mailtemplate.CodeSceneText(c.event.Scene)
		c.event.Content = rendered
		return "", rendered, nil
	}

	logger.Infof("使用邮件模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	if tmpl.Title != "" {
		c.event.Title = tmpl.Title
	}
	if tmpl.Content == "" {
		return tmpl.Code, c.event.Content, nil
	}

	title, intro := mailtemplate.CodeSceneText(c.event.Scene)
	if tmpl.Title != "" {
		title = tmpl.Title
	}
	params := make(map[string]string, len(c.event.Params)+2)
	for key, value := range c.event.Params {
		params[key] = value
	}
	params["title"] = title
	params["intro"] = intro

	rendered, renderErr := mailtemplate.Render(tmpl.Content, params)
	if renderErr != nil {
		logger.Errorf("渲染邮件模板失败: %v", renderErr)
		rendered, renderErr = mailtemplate.CodeEmail(c.event.Scene, c.event.Params["code"], c.event.Params["time"])
		if renderErr != nil {
			return "", "", renderErr
		}
		c.event.Title = title
	}
	c.event.Content = rendered

	return tmpl.Code, rendered, nil
}

// send 投递邮件；标题与正文在 resolve 阶段已写回事件。
func (c *emailChannel) send(_ context.Context, _, _ string) error {
	return c.svcCtx.EmailDeliver.DeliveryEmail(&emailx.EmailMessage{
		To:      []string{c.event.Email},
		CC:      []string{},
		Subject: c.event.Title,
		Content: c.event.Content,
	})
}

// isNonRetryableEmailError 判断是否为不可重试的错误。
// 判定依赖厂商文案：适配层改写错误信息会让它静默失效，退化为重试至死信。
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
