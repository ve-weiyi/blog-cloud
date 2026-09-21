package notificationservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateNotifyMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateNotifyMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateNotifyMessageLogic {
	return &UpdateNotifyMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新通知消息
func (l *UpdateNotifyMessageLogic) UpdateNotifyMessage(in *notificationrpc.UpdateNotifyMessageRequest) (*notificationrpc.UpdateNotifyMessageResponse, error) {
	fields := map[string]interface{}{
		"title":       in.Title,
		"content":     in.Content,
		"category":    in.Category,
		"level":       in.Level,
		"target_type": in.TargetType,
		"target_ids":  in.TargetIds,
		"updated_at":  time.Now(),
	}

	_, err := l.svcCtx.TNotifyMessageModel.UpdateFields(l.ctx, fields, "id = ? AND status = ?", in.Id, "draft")
	if err != nil {
		return nil, err
	}

	return &notificationrpc.UpdateNotifyMessageResponse{
		Success: true,
	}, nil
}
