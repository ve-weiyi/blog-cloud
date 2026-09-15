package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteLoginLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteLoginLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLoginLogsLogic {
	return &BatchDeleteLoginLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除用户登录日志
func (l *BatchDeleteLoginLogsLogic) BatchDeleteLoginLogs(in *syslogrpc.BatchDeleteLoginLogsRequest) (*syslogrpc.BatchDeleteLoginLogsResponse, error) {
	rows, err := l.svcCtx.TLoginLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchDeleteLoginLogsResponse{
		SuccessCount: rows,
	}, nil
}
