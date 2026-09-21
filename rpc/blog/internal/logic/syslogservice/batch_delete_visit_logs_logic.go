package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteVisitLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteVisitLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteVisitLogsLogic {
	return &BatchDeleteVisitLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除访问日志
func (l *BatchDeleteVisitLogsLogic) BatchDeleteVisitLogs(in *syslogrpc.BatchDeleteVisitLogsRequest) (*syslogrpc.BatchDeleteVisitLogsResponse, error) {
	rows, err := l.svcCtx.TVisitLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchDeleteVisitLogsResponse{
		SuccessCount: rows,
	}, nil
}
