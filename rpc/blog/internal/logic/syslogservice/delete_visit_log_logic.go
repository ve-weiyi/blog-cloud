package syslogservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteVisitLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVisitLogLogic {
	return &DeleteVisitLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除访问日志
func (l *DeleteVisitLogLogic) DeleteVisitLog(in *syslogrpc.DeleteVisitLogRequest) (*syslogrpc.DeleteVisitLogResponse, error) {
	rows, err := l.svcCtx.TVisitLogModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.DeleteVisitLogResponse{
		SuccessCount: rows,
	}, nil
}
