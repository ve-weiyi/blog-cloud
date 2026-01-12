package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesApiLogic {
	return &DeletesApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除接口
func (l *DeletesApiLogic) DeletesApi(in *permissionrpc.IdsReq) (*permissionrpc.BatchResp, error) {
	rows, err := l.svcCtx.TApiModel.Deletes(l.ctx, "id in (?) or parent_id in (?) ", in.Ids, in.Ids)
	if err != nil {
		return nil, err
	}

	return &permissionrpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
