package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDestroyAlbumsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDestroyAlbumsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDestroyAlbumsLogic {
	return &BatchDestroyAlbumsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除相册
func (l *BatchDestroyAlbumsLogic) BatchDestroyAlbums(in *mediarpc.BatchDestroyAlbumsRequest) (*mediarpc.BatchDestroyAlbumsResponse, error) {
	rows, err := l.svcCtx.TAlbumModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &mediarpc.BatchDestroyAlbumsResponse{SuccessCount: rows}, nil
}
