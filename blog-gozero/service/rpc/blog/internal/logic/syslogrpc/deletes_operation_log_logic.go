package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesOperationLogLogic {
	return &DeletesOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除操作记录
func (l *DeletesOperationLogLogic) DeletesOperationLog(in *syslogrpc.DeletesOperationLogReq) (*syslogrpc.DeletesOperationLogResp, error) {
	rows, err := l.svcCtx.TOperationLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.DeletesOperationLogResp{
		SuccessCount: rows,
	}, nil
}
