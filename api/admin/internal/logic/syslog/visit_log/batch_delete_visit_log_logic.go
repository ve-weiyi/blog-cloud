package visit_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type BatchDeleteVisitLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除访问日志
func NewBatchDeleteVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteVisitLogLogic {
	return &BatchDeleteVisitLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteVisitLogLogic) BatchDeleteVisitLog(req *types.DeleteVisitLogReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SyslogService.BatchDeleteVisitLogs(l.ctx, &syslogservice.BatchDeleteVisitLogsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
