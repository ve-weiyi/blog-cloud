package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type PublishNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发布统一通知消息
func NewPublishNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishNotifyMessageLogic {
	return &PublishNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PublishNotifyMessageLogic) PublishNotifyMessage(req *types.PublishNotifyMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.PublishNotifyMessage(l.ctx, &notificationservice.PublishNotifyMessageRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	var count int64
	if out.Success {
		count = 1
	}

	return &types.BatchResp{
		SuccessCount: count,
	}, nil
}
