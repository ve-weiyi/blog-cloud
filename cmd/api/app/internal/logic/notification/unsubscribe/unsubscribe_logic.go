package unsubscribe

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type UnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 执行退订
func NewUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsubscribeLogic {
	return &UnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Unsubscribe 执行退订：凭令牌，不需要登录。重复提交幂等。
func (l *UnsubscribeLogic) Unsubscribe(req *types.UnsubscribeReq) (resp *types.EmptyResp, err error) {
	if _, err := l.svcCtx.NotificationService.Unsubscribe(l.ctx,
		&notificationservice.UnsubscribeRequest{Token: req.Token}); err != nil {
		return nil, err
	}
	return &types.EmptyResp{}, nil
}
