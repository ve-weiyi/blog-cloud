package unsubscribe

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

type VerifyUnsubscribeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 校验退订令牌
func NewVerifyUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyUnsubscribeLogic {
	return &VerifyUnsubscribeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// VerifyUnsubscribe 供邮件确认页在用户点击前调用：只回答会退订什么，不改变状态。
func (l *VerifyUnsubscribeLogic) VerifyUnsubscribe(req *types.UnsubscribeVerifyReq) (resp *types.UnsubscribeVerifyResp, err error) {
	out, err := l.svcCtx.NotificationService.VerifyUnsubscribe(l.ctx,
		&notificationservice.VerifyUnsubscribeRequest{Token: req.Token})
	if err != nil {
		return nil, err
	}
	return &types.UnsubscribeVerifyResp{Group: out.Group, Channel: out.Channel}, nil
}
