package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type UpdateAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumLogic {
	return &UpdateAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新相册
func (l *UpdateAlbumLogic) UpdateAlbum(in *resourcerpc.UpdateAlbumRequest) (*resourcerpc.UpdateAlbumResponse, error) {
	entity, err := l.svcCtx.TAlbumModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.AlbumName = in.AlbumName
	entity.AlbumDesc = in.AlbumDesc
	entity.AlbumCover = in.AlbumCover
	entity.Status = in.Status

	_, err = l.svcCtx.TAlbumModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.UpdateAlbumResponse{Success: true}, nil
}
