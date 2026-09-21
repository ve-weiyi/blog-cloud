package menu

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type BatchDeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewBatchDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteMenuLogic {
	return &BatchDeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteMenuLogic) BatchDeleteMenu(req *types.DeleteMenuReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.AccessService.BatchDeleteMenus(l.ctx, &accessservice.BatchDeleteMenusRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
