package mqlogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

// ConsumeLoginLogic 登录日志消费者逻辑
// 监听登录事件，记录登录日志并更新用户最后登录信息
type ConsumeLoginLogic struct {
	svcCtx *svc.ServiceContext
}

// NewConsumeLoginLogic 创建登录日志消费者逻辑
func NewConsumeLoginLogic(svcCtx *svc.ServiceContext) *ConsumeLoginLogic {
	return &ConsumeLoginLogic{svcCtx: svcCtx}
}

// Consume 处理登录日志消息
func (l *ConsumeLoginLogic) Consume(ctx context.Context, event *mq.LoginEvent) error {
	logger := logx.WithContext(ctx)
	logger.Infof("收到登录日志消息: userId=%s, loginType=%s, status=%d",
		event.UserId, event.LoginType, event.Status)

	// 记录登录日志
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

	if _, err := l.svcCtx.TLoginLogModel.Insert(ctx, loginLog); err != nil {
		logger.Errorf("插入登录日志失败: %v", err)
		return err
	}

	return nil
}
