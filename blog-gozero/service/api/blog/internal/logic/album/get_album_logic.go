package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取相册
func NewGetAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlbumLogic {
	return &GetAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAlbumLogic) GetAlbum(req *types.IdReq) (resp *types.Album, err error) {
	in := &resourcerpc.GetAlbumReq{Id: req.Id}
	album, err := l.svcCtx.ResourceRpc.GetAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.Album{
		Id:         album.Album.Id,
		AlbumName:  album.Album.AlbumName,
		AlbumDesc:  album.Album.AlbumDesc,
		AlbumCover: album.Album.AlbumCover,
	}, nil
}
