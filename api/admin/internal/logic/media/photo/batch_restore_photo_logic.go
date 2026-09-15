package photo

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRestorePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量恢复照片
func NewBatchRestorePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRestorePhotoLogic {
	return &BatchRestorePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRestorePhotoLogic) BatchRestorePhoto(req *types.RestorePhotoReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchPatchPhotos(l.ctx, &mediaservice.BatchPatchPhotosRequest{
		Ids:      req.Ids,
		IsDelete: 0,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
