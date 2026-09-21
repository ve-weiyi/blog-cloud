package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteApisLogic {
	return &BatchDeleteApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除 API
func (l *BatchDeleteApisLogic) BatchDeleteApis(in *accessrpc.BatchDeleteApisRequest) (*accessrpc.BatchDeleteApisResponse, error) {
	if len(in.Ids) == 0 {
		return &accessrpc.BatchDeleteApisResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "api_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TApiModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &accessrpc.BatchDeleteApisResponse{SuccessCount: rows}, nil
}
