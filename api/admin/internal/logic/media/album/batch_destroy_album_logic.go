package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type BatchDestroyAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新相册删除状态
func NewBatchDestroyAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDestroyAlbumLogic {
	return &BatchDestroyAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDestroyAlbumLogic) BatchDestroyAlbum(req *types.DestroyAlbumReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchDestroyAlbums(l.ctx, &mediaservice.BatchDestroyAlbumsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
