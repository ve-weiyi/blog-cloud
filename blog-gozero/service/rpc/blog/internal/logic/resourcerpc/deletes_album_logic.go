package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletesAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesAlbumLogic {
	return &DeletesAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除相册
func (l *DeletesAlbumLogic) DeletesAlbum(in *resourcerpc.DeletesAlbumReq) (*resourcerpc.DeletesAlbumResp, error) {
	rows, err := l.svcCtx.TAlbumModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.TPhotoModel.Deletes(l.ctx, "album_id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.DeletesAlbumResp{
		SuccessCount: rows,
	}, nil
}
