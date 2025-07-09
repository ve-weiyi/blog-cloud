package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesVisitLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesVisitLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesVisitLogLogic {
	return &DeletesVisitLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除访问记录
func (l *DeletesVisitLogLogic) DeletesVisitLog(in *syslogrpc.IdsReq) (*syslogrpc.BatchResp, error) {
	rows, err := l.svcCtx.TVisitLogModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &syslogrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
