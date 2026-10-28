package consumer

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
)

// ConsumeLogoutLogic 登出消费者逻辑
// 监听登出事件，追加一条登出记录；登录记录只作为事件保留，不再回写。
type ConsumeLogoutLogic struct {
	svcCtx *svc.ServiceContext
}

// NewConsumeLogoutLogic 创建登出消费者逻辑
func NewConsumeLogoutLogic(svcCtx *svc.ServiceContext) *ConsumeLogoutLogic {
	return &ConsumeLogoutLogic{svcCtx: svcCtx}
}

// Consume 处理登出消息
func (l *ConsumeLogoutLogic) Consume(ctx context.Context, event *mq.LogoutEvent) error {
	logger := logx.WithContext(ctx)
	logger.Infof("收到登出消息: userId=%s, did=%s, logoutType=%s",
		event.UserId, event.DeviceId, event.LogoutType)

	now := time.Now()
	logout := &model.TLoginLog{
		UserId:     event.UserId,
		DeviceId:   event.DeviceId,
		EventType:  enums.AuthEventLogout,
		LogoutType: event.LogoutType,
		Status:     enums.AuthResultSuccess,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if _, err := l.svcCtx.TLoginLogModel.Insert(ctx, logout); err != nil {
		logger.Errorf("插入登出记录失败: %v", err)
		return err
	}

	return nil
}
