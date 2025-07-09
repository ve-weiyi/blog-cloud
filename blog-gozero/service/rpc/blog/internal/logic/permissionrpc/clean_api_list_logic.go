package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CleanApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCleanApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CleanApiListLogic {
	return &CleanApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 清空接口列表
func (l *CleanApiListLogic) CleanApiList(in *permissionrpc.CleanApiListReq) (*permissionrpc.BatchResp, error) {
	row, err := l.svcCtx.TApiModel.Deletes(l.ctx, "1 = 1")
	if err != nil {
		return nil, err
	}

	return &permissionrpc.BatchResp{
		SuccessCount: row,
	}, nil
}
