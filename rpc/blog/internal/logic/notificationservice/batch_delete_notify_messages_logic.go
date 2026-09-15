package notificationservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteNotifyMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteNotifyMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteNotifyMessagesLogic {
	return &BatchDeleteNotifyMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除通知消息
func (l *BatchDeleteNotifyMessagesLogic) BatchDeleteNotifyMessages(in *notificationrpc.BatchDeleteNotifyMessagesRequest) (*notificationrpc.BatchDeleteNotifyMessagesResponse, error) {
	if len(in.Ids) == 0 {
		return &notificationrpc.BatchDeleteNotifyMessagesResponse{}, nil
	}

	placeholders := make([]string, len(in.Ids))
	args := make([]interface{}, len(in.Ids))
	for i, id := range in.Ids {
		placeholders[i] = "?"
		args[i] = id
	}

	condition := fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ","))
	rows, err := l.svcCtx.TNotifyMessageModel.DeleteBatch(l.ctx, condition, args...)
	if err != nil {
		return nil, err
	}

	return &notificationrpc.BatchDeleteNotifyMessagesResponse{
		SuccessCount: rows,
	}, nil
}
