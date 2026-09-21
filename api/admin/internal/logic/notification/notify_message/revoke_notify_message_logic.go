package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type RevokeNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 撤回统一通知消息
func NewRevokeNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeNotifyMessageLogic {
	return &RevokeNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RevokeNotifyMessageLogic) RevokeNotifyMessage(req *types.RevokeNotifyMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.RevokeNotifyMessage(l.ctx, &notificationservice.RevokeNotifyMessageRequest{
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
