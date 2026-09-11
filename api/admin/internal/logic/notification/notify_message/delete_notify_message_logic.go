package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type DeleteNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除统一通知消息
func NewDeleteNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNotifyMessageLogic {
	return &DeleteNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNotifyMessageLogic) DeleteNotifyMessage(req *types.DeleteNotifyMessageReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.NotificationService.DeleteNotifyMessage(l.ctx, &notificationservice.DeleteNotifyMessageRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
