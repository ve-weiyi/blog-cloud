package mqlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// LoginConsumerLogic 登录日志消费者逻辑
// 监听登录事件，记录登录日志并更新用户最后登录信息
type LoginConsumerLogic struct {
	svcCtx   *svc.ServiceContext
	consumer mqx.Consumer
	logx.Logger
}

// NewLoginConsumerLogic 创建登录日志消费者逻辑
func NewLoginConsumerLogic(svcCtx *svc.ServiceContext) (*LoginConsumerLogic, error) {
	// 类型断言获取 MessageQueue
	loginMQ, ok := mq.LoginMQ.(mqx.MessageQueue)
	if !ok || loginMQ == nil {
		return nil, fmt.Errorf("LoginMQ 未初始化或类型不匹配")
	}

	// 创建消费者
	consumer, err := loginMQ.Consumer(&mqx.ConsumerConfig{
		Topics:        []string{mq.LoginQueue},
		ConsumerName:  "login-consumer",
		PrefetchCount: 10,
		MaxRetries:    3,
	})
	if err != nil {
		return nil, fmt.Errorf("创建登录日志消费者失败: %w", err)
	}

	return &LoginConsumerLogic{
		svcCtx:   svcCtx,
		consumer: consumer,
		Logger:   logx.WithContext(context.Background()),
	}, nil
}

// Start 启动消费者
func (l *LoginConsumerLogic) Start() {
	// 订阅主题
	err := l.consumer.Subscribe(mq.LoginRoutingKey)
	if err != nil {
		l.Logger.Errorf("订阅登录日志主题失败: %v", err)
		return
	}

	l.Logger.Info("登录日志消费者启动成功")

	// 使用处理器消费消息
	err = l.consumer.ConsumeWithHandler(mqx.MessageHandlerFunc(func(ctx context.Context, msg *mqx.Message) error {
		return l.handleMessage(msg.Body)
	}))

	if err != nil {
		l.Logger.Errorf("登录日志消费者运行失败: %v", err)
	}
}

// Stop 停止消费者
func (l *LoginConsumerLogic) Stop() {
	ctx := context.Background()
	if err := l.consumer.Stop(ctx); err != nil {
		l.Logger.Errorf("停止登录日志消费者失败: %v", err)
	}
}

// handleMessage 处理登录日志消息
func (l *LoginConsumerLogic) handleMessage(body []byte) error {
	// 1. 解析消息
	var event mq.LoginEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	l.Logger.Infof("收到登录日志消息: userId=%s, loginType=%s, status=%d",
		event.UserId, event.LoginType, event.Status)

	// 2. 创建登录日志记录
	ctx := context.Background()
	loginLog := &model.TLoginLog{
		Id:         0,
		UserId:     event.UserId,
		DeviceId:   event.DeviceId,
		LoginType:  event.LoginType,
		Status:     event.Status,
		FailReason: event.FailReason,
		LogoutAt:   nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 3. 插入数据库
	_, err := l.svcCtx.TLoginLogModel.Insert(ctx, loginLog)
	if err != nil {
		l.Logger.Errorf("插入登录日志失败: %v", err)
		return err
	}

	return nil
}
