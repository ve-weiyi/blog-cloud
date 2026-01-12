package login_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesLoginLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除登录日志
func NewDeletesLoginLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesLoginLogLogic {
	return &DeletesLoginLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesLoginLogLogic) DeletesLoginLog(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.DeletesLoginLogReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SyslogRpc.DeletesLoginLog(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
