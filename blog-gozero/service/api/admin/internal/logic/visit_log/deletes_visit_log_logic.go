package visit_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesVisitLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除操作记录
func NewDeletesVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesVisitLogLogic {
	return &DeletesVisitLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesVisitLogLogic) DeletesVisitLog(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SyslogRpc.DeletesVisitLog(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
