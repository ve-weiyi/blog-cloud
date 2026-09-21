package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type UpdateNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新统一通知消息
func NewUpdateNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyMessageLogic {
	return &UpdateNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateNotifyMessageLogic) UpdateNotifyMessage(req *types.UpdateNotifyMessageReq) (resp *types.NotifyMessageVO, err error) {
	_, err = l.svcCtx.NotificationService.UpdateNotifyMessage(l.ctx, &notificationservice.UpdateNotifyMessageRequest{
		Id:         req.Id,
		Title:      req.Title,
		Content:    req.Content,
		Category:   req.Category,
		Level:      req.Level,
		TargetType: req.TargetType,
		TargetIds:  req.TargetIds,
	})
	if err != nil {
		return nil, err
	}

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
