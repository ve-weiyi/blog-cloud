package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperationLogLogic {
	return &DeleteOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除操作记录
func (l *DeleteOperationLogLogic) DeleteOperationLog(in *syslogrpc.IdsReq) (*syslogrpc.BatchResp, error) {
	rows, err := l.svcCtx.TOperationLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
