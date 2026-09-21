package photo

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type BatchDeletePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除照片
func NewBatchDeletePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeletePhotoLogic {
	return &BatchDeletePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeletePhotoLogic) BatchDeletePhoto(req *types.DeletePhotoReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchPatchPhotos(l.ctx, &mediaservice.BatchPatchPhotosRequest{
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
