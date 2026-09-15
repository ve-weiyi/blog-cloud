package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteUploadLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteUploadLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteUploadLogsLogic {
	return &BatchDeleteUploadLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文件上传日志
func (l *BatchDeleteUploadLogsLogic) BatchDeleteUploadLogs(in *syslogrpc.BatchDeleteUploadLogsRequest) (*syslogrpc.BatchDeleteUploadLogsResponse, error) {
	rows, err := l.svcCtx.TUploadLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchDeleteUploadLogsResponse{
		SuccessCount: rows,
	}, nil
}
