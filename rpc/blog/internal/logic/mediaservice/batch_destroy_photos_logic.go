package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDestroyPhotosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDestroyPhotosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDestroyPhotosLogic {
	return &BatchDestroyPhotosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除照片
func (l *BatchDestroyPhotosLogic) BatchDestroyPhotos(in *mediarpc.BatchDestroyPhotosRequest) (*mediarpc.BatchDestroyPhotosResponse, error) {
	rows, err := l.svcCtx.TPhotoModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &mediarpc.BatchDestroyPhotosResponse{SuccessCount: rows}, nil
}
