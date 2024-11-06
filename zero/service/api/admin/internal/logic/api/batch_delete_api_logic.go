package api

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除api路由
func NewBatchDeleteApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteApiLogic {
	return &BatchDeleteApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteApiLogic) BatchDeleteApi(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &permissionrpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.PermissionRpc.DeleteApi(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
