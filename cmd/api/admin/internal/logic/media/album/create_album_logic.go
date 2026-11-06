package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/mediaservice"
)

type CreateAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建相册
func NewCreateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAlbumLogic {
	return &CreateAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateAlbumLogic) CreateAlbum(req *types.CreateAlbumReq) (resp *types.AlbumVO, err error) {
	out, err := l.svcCtx.MediaService.CreateAlbum(l.ctx, &mediaservice.CreateAlbumRequest{
		AlbumName:  req.AlbumName,
		AlbumDesc:  req.AlbumDesc,
		AlbumCover: req.AlbumCover,
		Status:     req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.AlbumVO{
		Id:         out.Album.Id,
		AlbumName:  out.Album.AlbumName,
		AlbumDesc:  out.Album.AlbumDesc,
		AlbumCover: out.Album.AlbumCover,
		Status:     out.Album.Status,
		PhotoCount: out.Album.PhotoCount,
		CreatedAt:  out.Album.CreatedAt,
		UpdatedAt:  out.Album.UpdatedAt,
	}, nil
}
