package userauthservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/mq"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/userauthrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/mqx"
	"github.com/ve-weiyi/vkit/x/jsonconv"

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
func (l *LogoutLogic) Logout(in *userauthrpc.LogoutRequest) (*userauthrpc.LogoutResponse, error) {
	uid, err := metax.GetUserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	did, err := metax.GetDeviceIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}

	// 推送退出登录消息
	if mq.LogoutProducer != nil {
		mq.LogoutProducer.Send(l.ctx, &mqx.Message{
			Topic: mq.LogoutQueue,
			Key:   mq.LogoutRoutingKey,
			Body: []byte(jsonconv.AnyToJsonNE(mq.LogoutEvent{
				UserId:     uid,
				DeviceId:   did,
				LogoutType: "user logout",
			})),
			Timestamp: time.Now(),
		})
	}

	return &userauthrpc.LogoutResponse{}, nil
}
