package notificationservicelogic

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
)

type VerifyUnsubscribeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyUnsubscribeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyUnsubscribeLogic {
	return &VerifyUnsubscribeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验退订令牌
//
// 邮件确认页在用户点击之前调用：只回答「会退订什么」，不改变任何状态。
func (l *VerifyUnsubscribeLogic) VerifyUnsubscribe(in *notificationrpc.VerifyUnsubscribeRequest) (*notificationrpc.VerifyUnsubscribeResponse, error) {
	claim, err := verifyToken(l.svcCtx.Config.NotifyConf.UnsubscribeSecret, in.Token)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.VerifyUnsubscribeResponse{
		Group:   string(claim.Group),
		Channel: string(claim.Channel),
	}, nil
}

// verifyToken 是校验与执行两条路径共用的入口。
func verifyToken(secret, token string) (notify.UnsubscribeClaim, error) {
	if secret == "" {
		return notify.UnsubscribeClaim{}, bizerr.NewBizError(bizcode.CodeServiceUnavailable, "退订功能未配置")
	}

	claim, err := notify.VerifyUnsubscribe(secret, token, time.Now())
	switch {
	case errors.Is(err, notify.ErrTokenExpired):
		return claim, bizerr.NewBizError(bizcode.CodeOperationNotAllowed, "退订链接已过期")
	case err != nil:
		return claim, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "退订链接无效")
	}
	return claim, nil
}
