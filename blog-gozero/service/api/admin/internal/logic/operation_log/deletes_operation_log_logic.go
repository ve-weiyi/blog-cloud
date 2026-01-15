package operation_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesOperationLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除操作记录
func NewDeletesOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesOperationLogLogic {
	return &DeletesOperationLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesOperationLogLogic) DeletesOperationLog(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.DeletesOperationLogReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SyslogRpc.DeletesOperationLog(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
