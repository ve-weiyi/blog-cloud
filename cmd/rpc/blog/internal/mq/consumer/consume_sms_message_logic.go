package consumer

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
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
	logx.WithContext(ctx).Infof("收到短信消息: mobile=%s, scene=%s", event.Mobile, event.Scene)

	return (&delivery{
		records: l.svcCtx.TNotifyRecordModel,
		ch:      &smsChannel{svcCtx: l.svcCtx, event: event},
		bizID:   eventx.EventIDFrom(ctx),
	}).run(ctx)
}

// smsChannel 把短信事件适配到统一投递流程。
type smsChannel struct {
	svcCtx *svc.ServiceContext
	event  *mq.SmsMessageEvent
}

func (c *smsChannel) name() string      { return enums.NotifyChannelSms }
func (c *smsChannel) recipient() string { return c.event.Mobile }

// retryable 把厂商文案判定的「不可重试」翻转成「是否值得重试」。
func (c *smsChannel) retryable(err error) bool { return !isNonRetryableSmsError(err) }

// resolve 优先使用库里的模板，缺失时回落到内置文案。
// 模板编码来自厂商侧配置，库里的非空编码优先。
func (c *smsChannel) resolve(ctx context.Context) (string, string, error) {
	logger := logx.WithContext(ctx)
	templateCode := c.svcCtx.SmsProvider.GetTemplateCode(c.event.Scene)

	tmpl, err := c.svcCtx.TNotifyTemplateModel.FindOne(ctx, "",
		"scene = ? AND channel = ? AND enabled = ?", c.event.Scene, enums.NotifyChannelSms, 1)
	if err != nil {
		logger.Infof("未找到短信模板(scene=%s)，使用默认内容", c.event.Scene)
		content := fmt.Sprintf("您的验证码是：%s，%s分钟内有效", c.event.Params["code"], c.event.Params["time"])
		return templateCode, content, nil
	}

	logger.Infof("使用短信模板: code=%s, scene=%s", tmpl.Code, tmpl.Scene)
	if tmpl.Code != "" {
		templateCode = tmpl.Code
	}
	return templateCode, tmpl.Content, nil
}

// send 投递短信。
func (c *smsChannel) send(ctx context.Context, templateCode, _ string) error {
	return c.svcCtx.SmsProvider.SendTemplate(ctx, c.event.Mobile, templateCode, c.event.Params)
}

// isNonRetryableSmsError 判断是否为不可重试的错误。
// 判定依赖厂商文案：适配层改写错误信息会让它静默失效，退化为重试至死信。
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
