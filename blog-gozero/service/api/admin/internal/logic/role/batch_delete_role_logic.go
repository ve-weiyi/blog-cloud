package role

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除角色
func NewBatchDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteRoleLogic {
	return &BatchDeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteRoleLogic) BatchDeleteRole(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &permissionrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.PermissionRpc.DeleteRole(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
