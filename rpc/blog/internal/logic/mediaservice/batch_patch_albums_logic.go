package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchPatchAlbumsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchPatchAlbumsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchPatchAlbumsLogic {
	return &BatchPatchAlbumsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 部分更新相册
func (l *BatchPatchAlbumsLogic) BatchPatchAlbums(in *mediarpc.BatchPatchAlbumsRequest) (*mediarpc.BatchPatchAlbumsResponse, error) {
	rows, err := l.svcCtx.TAlbumModel.UpdateFields(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.TPhotoModel.UpdateFields(l.ctx, map[string]interface{}{
		"is_delete": in.IsDelete,
	}, "album_id in (?)", in.Ids)
	if err != nil {
		l.Errorf("PatchAlbum TPhotoModel UpdateFields error: %v", err)
	}

	return &mediarpc.BatchPatchAlbumsResponse{SuccessCount: rows}, nil
}
