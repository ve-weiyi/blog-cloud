package consumer

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/mailtemplate"
	"github.com/ve-weiyi/vkit/adapter/emailx"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

// ConsumeEmailMessageLogic 邮件消息消费者：无状态发送，不写投递记录。
//
// 验证码的权威状态在 Redis，重发由用户触发；这里失败只需按可重试性交给 mqx。
type ConsumeEmailMessageLogic struct {
	svcCtx *svc.ServiceContext
}

func NewConsumeEmailMessageLogic(svcCtx *svc.ServiceContext) *ConsumeEmailMessageLogic {
	return &ConsumeEmailMessageLogic{svcCtx: svcCtx}
}

// Consume 处理邮件消息
func (l *ConsumeEmailMessageLogic) Consume(ctx context.Context, event *mq.EmailMessageEvent) error {
	logx.WithContext(ctx).Infof("收到邮件消息: email=%s, scene=%s", event.Email, event.Scene)

	ch := &emailChannel{svcCtx: l.svcCtx, event: event}
	if err := ch.resolve(ctx); err != nil {
		return err
	}
	if err := ch.send(ctx); err != nil {
		if !ch.retryable(err) {
			// 邮箱不存在、服务未开通这类错误重试只会进死信，直接确认
			return mqx.Permanent(err)
		}
		return err
	}
	return nil
}

// emailChannel 把邮件事件适配到发送流程。
type emailChannel struct {
	svcCtx *svc.ServiceContext
	event  *mq.EmailMessageEvent
}

// retryable 把厂商文案判定的「不可重试」翻转成「是否值得重试」。
func (c *emailChannel) retryable(err error) bool { return !isNonRetryableEmailError(err) }

// resolve 用代码里的模板渲染，模板缺失或渲染失败都回退到内置文案。
func (c *emailChannel) resolve(ctx context.Context) error {
	rendered, err := mailtemplate.CodeEmail(c.event.Scene, c.event.Params["code"], c.event.Params["time"])
	if err != nil {
		return err
	}
	title, _ := mailtemplate.CodeSceneText(c.event.Scene)
	c.event.Title = title
	c.event.Content = rendered
	return nil
}

// send 投递邮件；标题与正文在 resolve 阶段已写回事件。
func (c *emailChannel) send(_ context.Context) error {
	return c.svcCtx.EmailDeliver.DeliveryEmail(&emailx.EmailMessage{
		To:      []string{c.event.Email},
		CC:      []string{},
		Subject: c.event.Title,
		Content: c.event.Content,
	})
}

// isNonRetryableEmailError 判断是否为不可重试的错误。
// 依赖厂商文案：适配层改写错误信息会让它静默失效，退化为重试至死信。
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
