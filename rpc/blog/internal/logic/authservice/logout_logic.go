package authservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/authrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"

	"github.com/ve-weiyi/blog-cloud/infra/metax"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 退出登录
func (l *LogoutLogic) Logout(in *authrpc.LogoutRequest) (*authrpc.LogoutResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	did, err := metax.GetDeviceIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 推送退出登录消息，失败不影响登出流程
	if err := mq.PublishLogoutEvent(l.ctx, &mq.LogoutEvent{
		UserId:     uid,
		DeviceId:   did,
		LogoutType: "user logout",
	}); err != nil {
		l.Logger.Errorf("发送登出消息失败: %v", err)
	}

	return &authrpc.LogoutResponse{}, nil
}
