package accessservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/accessrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteMenusLogic {
	return &BatchDeleteMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除菜单
func (l *BatchDeleteMenusLogic) BatchDeleteMenus(in *accessrpc.BatchDeleteMenusRequest) (*accessrpc.BatchDeleteMenusResponse, error) {
	if len(in.Ids) == 0 {
		return &accessrpc.BatchDeleteMenusResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "menu_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TMenuModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &accessrpc.BatchDeleteMenusResponse{SuccessCount: rows}, nil
}
