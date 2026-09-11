package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type PatchAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPatchAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PatchAlbumLogic {
	return &PatchAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 部分更新相册
func (l *PatchAlbumLogic) PatchAlbum(in *resourcerpc.PatchAlbumRequest) (*resourcerpc.PatchAlbumResponse, error) {
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

	return &resourcerpc.PatchAlbumResponse{SuccessCount: rows}, nil
}
