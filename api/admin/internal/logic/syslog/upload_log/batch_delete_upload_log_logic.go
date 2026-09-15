package upload_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type BatchDeleteUploadLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除文件日志
func NewBatchDeleteUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteUploadLogLogic {
	return &BatchDeleteUploadLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteUploadLogLogic) BatchDeleteUploadLog(req *types.DeleteUploadLogReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SyslogService.BatchDeleteUploadLogs(l.ctx, &syslogservice.BatchDeleteUploadLogsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
