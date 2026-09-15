package operation_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type BatchDeleteOperationLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除操作日志
func NewBatchDeleteOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperationLogLogic {
	return &BatchDeleteOperationLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteOperationLogLogic) BatchDeleteOperationLog(req *types.DeleteOperationLogReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SyslogService.BatchDeleteOperationLogs(l.ctx, &syslogservice.BatchDeleteOperationLogsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
