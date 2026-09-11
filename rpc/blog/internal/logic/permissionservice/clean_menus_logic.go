package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type CleanMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanMenusLogic {
	return &CleanMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空菜单
func (l *CleanMenusLogic) CleanMenus(in *permissionrpc.CleanMenusRequest) (*permissionrpc.CleanMenusResponse, error) {
	_, err := l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TMenuModel.DeleteBatch(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	return &permissionrpc.CleanMenusResponse{SuccessCount: rows}, nil
}
