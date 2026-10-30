package consumer

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

// ConsumeSmsMessageLogic 短信消息消费者：无状态发送，不写投递记录。
//
// 模板号由厂商侧配置提供，正文只用于日志，实际内容由服务商按模板渲染。
type ConsumeSmsMessageLogic struct {
	svcCtx *svc.ServiceContext
}

func NewConsumeSmsMessageLogic(svcCtx *svc.ServiceContext) *ConsumeSmsMessageLogic {
	return &ConsumeSmsMessageLogic{svcCtx: svcCtx}
}

// Consume 处理短信消息
func (l *ConsumeSmsMessageLogic) Consume(ctx context.Context, event *mq.SmsMessageEvent) error {
	logx.WithContext(ctx).Infof("收到短信消息: mobile=%s, scene=%s", event.Mobile, event.Scene)

	code := l.svcCtx.SmsProvider.GetTemplateCode(event.Scene)
	err := l.svcCtx.SmsProvider.SendTemplate(ctx, event.Mobile, code, event.Params)
	if err != nil {
		if isNonRetryableSmsError(err) {
			// 号码非法、模板不存在这类错误重试只会进死信，直接确认
			return mqx.Permanent(err)
		}
		return err
	}
	return nil
}

// isNonRetryableSmsError 判断是否为不可重试的错误。
// 依赖厂商文案：适配层改写错误信息会让它静默失效，退化为重试至死信。
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
