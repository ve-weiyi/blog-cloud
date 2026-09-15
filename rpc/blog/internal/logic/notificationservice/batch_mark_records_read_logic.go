package notificationservicelogic

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchMarkRecordsReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchMarkRecordsReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchMarkRecordsReadLogic {
	return &BatchMarkRecordsReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchMarkRecordsReadLogic) BatchMarkRecordsRead(in *notificationrpc.BatchMarkRecordsReadRequest) (*notificationrpc.BatchMarkRecordsReadResponse, error) {
	if len(in.Ids) == 0 {
		return &notificationrpc.BatchMarkRecordsReadResponse{}, nil
	}

	placeholders := make([]string, len(in.Ids))
	args := make([]any, 0, len(in.Ids)+1)
	args = append(args, "unread")
	for i, id := range in.Ids {
		placeholders[i] = "?"
		args = append(args, id)
	}

	fields := map[string]interface{}{
		"status":  "read",
		"read_at": time.Now(),
	}
	rows, err := l.svcCtx.TNotifyRecordModel.UpdateFields(l.ctx, fields,
		"status = ? AND id IN ("+strings.Join(placeholders, ",")+")", args...)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.BatchMarkRecordsReadResponse{
		SuccessCount: rows,
	}, nil
}
