package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除相册
func NewBatchDeleteAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteAlbumLogic {
	return &BatchDeleteAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteAlbumLogic) BatchDeleteAlbum(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &photorpc.IdsReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.PhotoRpc.DeleteAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
