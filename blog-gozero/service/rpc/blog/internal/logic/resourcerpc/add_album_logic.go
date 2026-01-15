package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAlbumLogic {
	return &AddAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建相册
func (l *AddAlbumLogic) AddAlbum(in *resourcerpc.AddAlbumReq) (*resourcerpc.AddAlbumResp, error) {
	entity := &model.TAlbum{
		Id:         in.Id,
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		IsDelete:   in.IsDelete,
		Status:     in.Status,
	}

	_, err := l.svcCtx.TAlbumModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.AddAlbumResp{
		Album: convertAlbumOut(entity, nil),
	}, nil
}

func convertAlbumOut(in *model.TAlbum, cm map[int64]int) (out *resourcerpc.Album) {
	var count int
	if v, ok := cm[in.Id]; ok {
		count = v
	}

	out = &resourcerpc.Album{
		Id:         in.Id,
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		IsDelete:   in.IsDelete,
		Status:     in.Status,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
		PhotoCount: int64(count),
	}

	return out
}
