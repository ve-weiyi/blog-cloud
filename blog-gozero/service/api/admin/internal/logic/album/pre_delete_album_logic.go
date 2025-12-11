package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreDeleteAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 预删除相册
func NewPreDeleteAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreDeleteAlbumLogic {
	return &PreDeleteAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreDeleteAlbumLogic) PreDeleteAlbum(req *types.PreDeleteAlbumReq) (resp *types.BatchResp, err error) {
	in := &resourcerpc.UpdateAlbumDeleteReq{
		Ids:      req.Ids,
		IsDelete: req.IsDelete,
	}
	out, err := l.svcCtx.ResourceRpc.UpdateAlbumDelete(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
