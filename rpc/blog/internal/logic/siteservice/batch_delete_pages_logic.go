package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeletePagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeletePagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeletePagesLogic {
	return &BatchDeletePagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除页面
func (l *BatchDeletePagesLogic) BatchDeletePages(in *siterpc.BatchDeletePagesRequest) (*siterpc.BatchDeletePagesResponse, error) {
	rows, err := l.svcCtx.TPageModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &siterpc.BatchDeletePagesResponse{SuccessCount: rows}, nil
}
