package operation_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteOperationLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除操作记录
func NewBatchDeleteOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteOperationLogLogic {
	return &BatchDeleteOperationLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteOperationLogLogic) BatchDeleteOperationLog(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SyslogRpc.DeleteOperationLog(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
