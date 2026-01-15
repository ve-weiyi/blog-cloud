package file_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesFileLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除文件日志
func NewDeletesFileLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesFileLogLogic {
	return &DeletesFileLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesFileLogLogic) DeletesFileLog(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &syslogrpc.DeletesFileLogReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.SyslogRpc.DeletesFileLog(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
