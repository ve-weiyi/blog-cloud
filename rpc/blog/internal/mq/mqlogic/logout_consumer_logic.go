package mqlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx"
)

// LogoutConsumerLogic 登出消费者逻辑
// 监听登出事件，更新登录日志的登出时间
type LogoutConsumerLogic struct {
	svcCtx   *svc.ServiceContext
	consumer mqx.Consumer
	logx.Logger
}

// NewLogoutConsumerLogic 创建登出消费者逻辑
func NewLogoutConsumerLogic(svcCtx *svc.ServiceContext) (*LogoutConsumerLogic, error) {
	// 类型断言获取 MessageQueue
	logoutMQ, ok := mq.LogoutMQ.(mqx.MessageQueue)
	if !ok || logoutMQ == nil {
		return nil, fmt.Errorf("LogoutMQ 未初始化或类型不匹配")
	}

	// 创建消费者
	consumer, err := logoutMQ.Consumer(&mqx.ConsumerConfig{
		Topics:        []string{mq.LogoutQueue},
		ConsumerName:  "logout-consumer",
		PrefetchCount: 10,
		MaxRetries:    3,
	})
	if err != nil {
		return nil, fmt.Errorf("创建登出消费者失败: %w", err)
	}

	return &LogoutConsumerLogic{
		svcCtx:   svcCtx,
		consumer: consumer,
		Logger:   logx.WithContext(context.Background()),
	}, nil
}

// Start 启动消费者
func (l *LogoutConsumerLogic) Start() {
	// 订阅主题
	err := l.consumer.Subscribe(mq.LogoutRoutingKey)
	if err != nil {
		l.Logger.Errorf("订阅登出主题失败: %v", err)
		return
	}

	l.Logger.Info("登出消费者启动成功")

	// 使用处理器消费消息
	err = l.consumer.ConsumeWithHandler(mqx.MessageHandlerFunc(func(ctx context.Context, msg *mqx.Message) error {
		return l.handleMessage(msg.Body)
	}))

	if err != nil {
		l.Logger.Errorf("登出消费者运行失败: %v", err)
	}
}

// Stop 停止消费者
func (l *LogoutConsumerLogic) Stop() {
	ctx := context.Background()
	if err := l.consumer.Stop(ctx); err != nil {
		l.Logger.Errorf("停止登出消费者失败: %v", err)
	}
}

// handleMessage 处理登出消息
func (l *LogoutConsumerLogic) handleMessage(body []byte) error {
	// 1. 解析消息
	var event mq.LogoutEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("解析消息失败: %v", err)
	}

	l.Logger.Infof("收到登出消息: userId=%s, did=%s, logoutType=%s",
		event.UserId, event.DeviceId, event.LogoutType)

	// 2. 查找最近的登录记录（未登出的）
	ctx := context.Background()
	exists, _, err := l.svcCtx.TLoginLogModel.FindListAndTotal(ctx, 1, 1, "id desc", "user_id = ?", event.UserId)
	if err != nil {
		l.Logger.Errorf("查找登录记录失败: userId=%s, error=%v", event.UserId, err)
		// 如果找不到登录记录，不返回错误，只记录日志
		return nil
	}
	if len(exists) == 0 {
		l.Logger.Infof("未找到登录记录: userId=%s", event.UserId)
		return nil
	}
	loginLog := exists[0]

	// 3. 更新登出时间
	now := time.Now()
	loginLog.LogoutAt = &now
	loginLog.UpdatedAt = now

	_, err = l.svcCtx.TLoginLogModel.Update(ctx, loginLog)
	if err != nil {
		l.Logger.Errorf("更新登出时间失败: %v", err)
		return err
	}

	l.Logger.Infof("用户登出成功: userId=%s, did=%s, logoutType=%s",
		event.UserId, event.DeviceId, event.LogoutType)
	return nil
}
