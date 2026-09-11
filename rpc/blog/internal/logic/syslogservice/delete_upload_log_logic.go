package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteUploadLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUploadLogLogic {
	return &DeleteUploadLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除文件上传日志
func (l *DeleteUploadLogLogic) DeleteUploadLog(in *syslogrpc.DeleteUploadLogRequest) (*syslogrpc.DeleteUploadLogResponse, error) {
	rows, err := l.svcCtx.TUploadLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.DeleteUploadLogResponse{
		SuccessCount: rows,
	}, nil
}
