package menu

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除菜单
func NewBatchDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteMenuLogic {
	return &BatchDeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteMenuLogic) BatchDeleteMenu(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &permissionrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.PermissionRpc.DeleteMenu(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
