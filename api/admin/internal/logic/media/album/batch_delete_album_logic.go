package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
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

func (l *BatchDeleteAlbumLogic) BatchDeleteAlbum(req *types.DeleteAlbumReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchPatchAlbums(l.ctx, &mediaservice.BatchPatchAlbumsRequest{
		Ids:      req.Ids,
		IsDelete: 1,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
