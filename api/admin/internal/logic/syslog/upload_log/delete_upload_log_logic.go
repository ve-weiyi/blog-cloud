package upload_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type DeleteUploadLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除文件日志
func NewDeleteUploadLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUploadLogLogic {
	return &DeleteUploadLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUploadLogLogic) DeleteUploadLog(req *types.DeleteUploadLogReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.SyslogService.DeleteUploadLog(l.ctx, &syslogservice.DeleteUploadLogRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
