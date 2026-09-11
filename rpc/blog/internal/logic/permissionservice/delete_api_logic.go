package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteApiLogic {
	return &DeleteApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除 API
func (l *DeleteApiLogic) DeleteApi(in *permissionrpc.DeleteApiRequest) (*permissionrpc.DeleteApiResponse, error) {
	if len(in.Ids) == 0 {
		return &permissionrpc.DeleteApiResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "api_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TApiModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.DeleteApiResponse{SuccessCount: rows}, nil
}
