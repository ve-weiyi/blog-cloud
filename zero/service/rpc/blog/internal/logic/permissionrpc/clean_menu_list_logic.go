package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanMenuListLogic {
	return &CleanMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空菜单列表
func (l *CleanMenuListLogic) CleanMenuList(in *permissionrpc.EmptyReq) (*permissionrpc.BatchResp, error) {
	row, err := l.svcCtx.TMenuModel.Deletes(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	return &permissionrpc.BatchResp{
		SuccessCount: row,
	}, nil
}
