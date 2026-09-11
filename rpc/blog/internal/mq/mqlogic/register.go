package mqlogic

import (
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

// RegisterMqConsumers 注册所有消息队列消费者
func RegisterMqConsumers(svcCtx *svc.ServiceContext) {
	logx.Info("Registering MQ consumers...")

	consumers := []struct {
		name string
		fn   func(*svc.ServiceContext)
	}{
		{"Email", func(svcCtx *svc.ServiceContext) {
			logic, err := NewEmailConsumerLogic(svcCtx)
			if err != nil {
				logx.Errorf("创建邮件消费者失败: %v", err)
				return
			}
			logic.Start()
		}},
		{"SMS", func(svcCtx *svc.ServiceContext) {
			logic, err := NewSmsConsumerLogic(svcCtx)
			if err != nil {
				logx.Errorf("创建短信消费者失败: %v", err)
				return
			}
			logic.Start()
		}},
		{"Login", func(svcCtx *svc.ServiceContext) {
			logic, err := NewLoginConsumerLogic(svcCtx)
			if err != nil {
				logx.Errorf("创建登录日志消费者失败: %v", err)
				return
			}
			logic.Start()
		}},
		{"Logout", func(svcCtx *svc.ServiceContext) {
			logic, err := NewLogoutConsumerLogic(svcCtx)
			if err != nil {
				logx.Errorf("创建登出消费者失败: %v", err)
				return
			}
			logic.Start()
		}},
		{"Inbox", func(svcCtx *svc.ServiceContext) {
			logic, err := NewInboxConsumerLogic(svcCtx)
			if err != nil {
				logx.Errorf("创建站内信消费者失败: %v", err)
				return
			}
			logic.Start()
		}},
	}

	for _, consumer := range consumers {
		go func(name string, fn func(*svc.ServiceContext)) {
			logx.Infof("Starting %s consumer...", name)
			fn(svcCtx)
		}(consumer.name, consumer.fn)
	}

	logx.Info("All MQ consumers registered successfully")
}
