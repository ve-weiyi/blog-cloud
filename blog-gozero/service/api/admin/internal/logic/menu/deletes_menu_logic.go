package menu

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单
func NewDeletesMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesMenuLogic {
	return &DeletesMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesMenuLogic) DeletesMenu(req *types.IdsReq) (resp *types.BatchResp, err error) {
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
