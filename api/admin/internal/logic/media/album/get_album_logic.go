package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
)

type GetAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取相册详情
func NewGetAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlbumLogic {
	return &GetAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAlbumLogic) GetAlbum(req *types.GetAlbumReq) (resp *types.AlbumVO, err error) {
	out, err := l.svcCtx.MediaService.GetAlbum(l.ctx, &mediaservice.GetAlbumRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.AlbumVO{
		Id:         out.Album.Id,
		AlbumName:  out.Album.AlbumName,
		AlbumDesc:  out.Album.AlbumDesc,
		AlbumCover: out.Album.AlbumCover,
		IsDelete:   out.Album.IsDelete,
		Status:     out.Album.Status,
		PhotoCount: out.Album.PhotoCount,
		CreatedAt:  out.Album.CreatedAt,
		UpdatedAt:  out.Album.UpdatedAt,
	}, nil
}
