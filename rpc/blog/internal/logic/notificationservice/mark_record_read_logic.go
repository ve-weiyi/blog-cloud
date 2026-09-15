package notificationservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type MarkRecordReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkRecordReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkRecordReadLogic {
	return &MarkRecordReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记单条 inbox 已读（按投递记录 id）
func (l *MarkRecordReadLogic) MarkRecordRead(in *notificationrpc.MarkRecordReadRequest) (*notificationrpc.MarkRecordReadResponse, error) {
	fields := map[string]interface{}{
		"status":  "read",
		"read_at": time.Now(),
	}

	_, err := l.svcCtx.TNotifyRecordModel.UpdateFields(l.ctx, fields,
		"id = ? AND status = ?",
		in.Id, "unread")
	if err != nil {
		return nil, err
	}

	return &notificationrpc.MarkRecordReadResponse{
		Success: true,
	}, nil
}
