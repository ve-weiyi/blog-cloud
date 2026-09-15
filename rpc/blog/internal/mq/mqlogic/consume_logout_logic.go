package mqlogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

// ConsumeLogoutLogic 登出消费者逻辑
// 监听登出事件，更新登录日志的登出时间
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

	// 查找最近的登录记录（未登出的）
	exists, _, err := l.svcCtx.TLoginLogModel.FindListAndTotal(ctx, 1, 1, "id desc", "user_id = ?", event.UserId)
	if err != nil {
		// 如果找不到登录记录，不返回错误，只记录日志
		logger.Errorf("查找登录记录失败: userId=%s, error=%v", event.UserId, err)
		return nil
	}
	if len(exists) == 0 {
		logger.Infof("未找到登录记录: userId=%s", event.UserId)
		return nil
	}
	loginLog := exists[0]

	// 更新登出时间
	now := time.Now()
	loginLog.LogoutAt = &now
	loginLog.UpdatedAt = now

	_, err = l.svcCtx.TLoginLogModel.Update(ctx, loginLog)
	if err != nil {
		logger.Errorf("更新登出时间失败: %v", err)
		return err
	}

	logger.Infof("用户登出成功: userId=%s, did=%s, logoutType=%s",
		event.UserId, event.DeviceId, event.LogoutType)
	return nil
}
