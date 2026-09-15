package album

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchRestoreAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量恢复相册
func NewBatchRestoreAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchRestoreAlbumLogic {
	return &BatchRestoreAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchRestoreAlbumLogic) BatchRestoreAlbum(req *types.RestoreAlbumReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.MediaService.BatchPatchAlbums(l.ctx, &mediaservice.BatchPatchAlbumsRequest{
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
