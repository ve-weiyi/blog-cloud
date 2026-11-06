package notificationservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
)

type UnsubscribeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsubscribeLogic {
	return &UnsubscribeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 凭令牌执行退订
//
// 不要求登录：令牌自带到收件人、分组与渠道，签名保证别人伪造不了。
// 幂等，重复提交不报错。
func (l *UnsubscribeLogic) Unsubscribe(in *notificationrpc.UnsubscribeRequest) (*notificationrpc.UnsubscribeResponse, error) {
	claim, err := verifyToken(l.svcCtx.Config.NotifyConf.UnsubscribeSecret, in.Token)
	if err != nil {
		return nil, err
	}

	if _, err := l.svcCtx.Notify.UpdatePreferences(l.ctx, claim.UserID, []notification.Preference{{
		Group:   claim.Group,
		Channel: claim.Channel,
		Enabled: false,
	}}); err != nil {
		return nil, err
	}

	return &notificationrpc.UnsubscribeResponse{
		Success: true,
		Group:   string(claim.Group),
		Channel: string(claim.Channel),
	}, nil
}
