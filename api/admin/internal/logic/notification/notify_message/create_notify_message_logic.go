package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/notificationservice"
)

type CreateNotifyMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建统一通知消息
func NewCreateNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNotifyMessageLogic {
	return &CreateNotifyMessageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateNotifyMessageLogic) CreateNotifyMessage(req *types.CreateNotifyMessageReq) (resp *types.NotifyMessageVO, err error) {
	out, err := l.svcCtx.NotificationService.CreateNotifyMessage(l.ctx, &notificationservice.CreateNotifyMessageRequest{
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

	detail, err := l.svcCtx.NotificationService.GetNotifyMessage(l.ctx, &notificationservice.GetNotifyMessageRequest{
		Id: out.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.NotifyMessageVO{
		Id:          detail.Message.Id,
		Title:       detail.Message.Title,
		Content:     detail.Message.Content,
		Category:    detail.Message.Category,
		Level:       detail.Message.Level,
		TargetType:  detail.Message.TargetType,
		TargetIds:   detail.Message.TargetIds,
		Status:      detail.Message.Status,
		PublishedAt: detail.Message.PublishedAt,
		PublishedBy: detail.Message.PublishedBy,
		CreatedAt:   detail.Message.CreatedAt,
		UpdatedAt:   detail.Message.UpdatedAt,
	}, nil
}
