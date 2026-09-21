package login_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type BatchDeleteLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除登录日志
func NewBatchDeleteLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteLoginLogLogic {
	return &BatchDeleteLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteLoginLogLogic) BatchDeleteLoginLog(req *types.DeleteLoginLogReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SyslogService.BatchDeleteLoginLogs(l.ctx, &syslogservice.BatchDeleteLoginLogsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
