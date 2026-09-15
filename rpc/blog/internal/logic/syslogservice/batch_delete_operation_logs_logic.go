package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteOperationLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteOperationLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperationLogsLogic {
	return &BatchDeleteOperationLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除操作日志
func (l *BatchDeleteOperationLogsLogic) BatchDeleteOperationLogs(in *syslogrpc.BatchDeleteOperationLogsRequest) (*syslogrpc.BatchDeleteOperationLogsResponse, error) {
	rows, err := l.svcCtx.TOperationLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchDeleteOperationLogsResponse{
		SuccessCount: rows,
	}, nil
}
