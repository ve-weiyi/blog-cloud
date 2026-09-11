package mq

import (
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// 全局 MQ 实例
var (
	EmailMQ  mqx.MessageQueue
	SmsMQ    mqx.MessageQueue
	InboxMQ  mqx.MessageQueue
	LoginMQ  mqx.MessageQueue
	LogoutMQ mqx.MessageQueue
)

// 全局 Producer 实例（复用，不每次创建）
var (
	EmailProducer  mqx.Producer
	SmsProducer    mqx.Producer
	InboxProducer  mqx.Producer
	LoginProducer  mqx.Producer
	LogoutProducer mqx.Producer
)

func Init(svcCtx *svc.ServiceContext) {
	var err error

	EmailMQ, err = initEmailMq(svcCtx)
	if err != nil {
		logx.Errorf("init email mq error: %s", err.Error())
	} else {
		EmailProducer, err = EmailMQ.Producer()
		if err != nil {
			logx.Errorf("init email producer error: %s", err.Error())
		}
	}

	SmsMQ, err = initSmsMq(svcCtx)
	if err != nil {
		logx.Errorf("init sms mq error: %s", err.Error())
	} else {
		SmsProducer, err = SmsMQ.Producer()
		if err != nil {
			logx.Errorf("init sms producer error: %s", err.Error())
		}
	}

	InboxMQ, err = initInboxMq(svcCtx)
	if err != nil {
		logx.Errorf("init inbox mq error: %s", err.Error())
	} else {
		InboxProducer, err = InboxMQ.Producer()
		if err != nil {
			logx.Errorf("init inbox producer error: %s", err.Error())
		}
	}

	LoginMQ, err = initLoginMq(svcCtx)
	if err != nil {
		logx.Errorf("init login mq error: %s", err.Error())
	} else {
		LoginProducer, err = LoginMQ.Producer()
		if err != nil {
			logx.Errorf("init login producer error: %s", err.Error())
		}
	}

	LogoutMQ, err = initLogoutMq(svcCtx)
	if err != nil {
		logx.Errorf("init logout mq error: %s", err.Error())
	} else {
		LogoutProducer, err = LogoutMQ.Producer()
		if err != nil {
			logx.Errorf("init logout producer error: %s", err.Error())
		}
	}
}
