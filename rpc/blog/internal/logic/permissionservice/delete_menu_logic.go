package permissionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除菜单
func (l *DeleteMenuLogic) DeleteMenu(in *permissionrpc.DeleteMenuRequest) (*permissionrpc.DeleteMenuResponse, error) {
	if len(in.Ids) == 0 {
		return &permissionrpc.DeleteMenuResponse{SuccessCount: 0}, nil
	}

	_, err := l.svcCtx.TRoleMenuModel.DeleteBatch(l.ctx, "menu_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	rows, err := l.svcCtx.TMenuModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.DeleteMenuResponse{SuccessCount: rows}, nil
}
