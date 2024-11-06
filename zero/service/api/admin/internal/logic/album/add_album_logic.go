package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建相册
func NewAddAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAlbumLogic {
	return &AddAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAlbumLogic) AddAlbum(req *types.AlbumNewReq) (resp *types.AlbumBackDTO, err error) {
	in := ConvertAlbumPb(req)
	out, err := l.svcCtx.PhotoRpc.AddAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertAlbumTypes(out)
	return resp, nil
}

func ConvertAlbumPb(in *types.AlbumNewReq) (out *photorpc.AlbumNewReq) {
	out = &photorpc.AlbumNewReq{
		Id:         in.Id,
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		IsDelete:   in.IsDelete,
		Status:     in.Status,
	}

	return
}

func ConvertAlbumTypes(in *photorpc.AlbumDetails) (out *types.AlbumBackDTO) {
	out = &types.AlbumBackDTO{
		Id:         in.Id,
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		IsDelete:   in.IsDelete,
		Status:     in.Status,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
		PhotoCount: in.PhotoCount,
	}

	return
}
