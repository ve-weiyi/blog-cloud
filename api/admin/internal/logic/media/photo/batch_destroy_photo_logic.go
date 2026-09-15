package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type BatchDestroyPhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量更新照片删除状态
func NewBatchDestroyPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDestroyPhotoLogic {
	return &BatchDestroyPhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDestroyPhotoLogic) BatchDestroyPhoto(req *types.DestroyPhotoReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchDestroyPhotos(l.ctx, &mediaservice.BatchDestroyPhotosRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
