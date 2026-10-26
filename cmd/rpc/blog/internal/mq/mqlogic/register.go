package mqlogic

import (
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx/event"
)

// Register 把全部事件处理器注册到 Router。
// 任一注册失败都返回错误：缺处理器的 routing key 会按失败策略处理，宁可不起订阅。
func Register(router *event.Router, svcCtx *svc.ServiceContext) error {
	for _, err := range []error{
		event.HandleEvent(router, mq.EmailMessage, NewConsumeEmailMessageLogic(svcCtx).Consume),
		event.HandleEvent(router, mq.SmsMessage, NewConsumeSmsMessageLogic(svcCtx).Consume),
		event.HandleEvent(router, mq.InboxMessage, NewConsumeInboxMessageLogic(svcCtx).Consume),
		event.HandleEvent(router, mq.LoginMessage, NewConsumeLoginLogic(svcCtx).Consume),
		event.HandleEvent(router, mq.LogoutMessage, NewConsumeLogoutLogic(svcCtx).Consume),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}
