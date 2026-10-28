package notify_message

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
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
