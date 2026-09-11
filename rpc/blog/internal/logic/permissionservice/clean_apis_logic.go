package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CleanApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanApisLogic {
	return &CleanApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空 API
func (l *CleanApisLogic) CleanApis(in *permissionrpc.CleanApisRequest) (*permissionrpc.CleanApisResponse, error) {
	_, err := l.svcCtx.TRoleApiModel.DeleteBatch(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TApiModel.DeleteBatch(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	return &permissionrpc.CleanApisResponse{SuccessCount: rows}, nil
}
