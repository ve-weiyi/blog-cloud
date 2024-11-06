package operation_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperationLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除操作记录
func NewDeleteOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperationLogLogic {
	return &DeleteOperationLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteOperationLogLogic) DeleteOperationLog(req *types.IdReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.IdsReq{
		Ids: []int64{req.Id},
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
