package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type GetNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取统一通知消息详情
func NewGetNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNotifyMessageLogic {
	return &GetNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNotifyMessageLogic) GetNotifyMessage(req *types.GetNotifyMessageReq) (resp *types.NotifyMessageVO, err error) {
	out, err := l.svcCtx.NotificationService.GetNotifyMessage(l.ctx, &notificationservice.GetNotifyMessageRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.NotifyMessageVO{
		Id:          out.Message.Id,
		Title:       out.Message.Title,
		Content:     out.Message.Content,
		Category:    out.Message.Category,
		Level:       out.Message.Level,
		TargetType:  out.Message.TargetType,
		TargetIds:   out.Message.TargetIds,
		Status:      out.Message.Status,
		PublishedAt: out.Message.PublishedAt,
		PublishedBy: out.Message.PublishedBy,
		CreatedAt:   out.Message.CreatedAt,
		UpdatedAt:   out.Message.UpdatedAt,
	}, nil
}
