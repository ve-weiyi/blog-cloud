package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询相册
func NewGetAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlbumLogic {
	return &GetAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAlbumLogic) GetAlbum(req *types.IdReq) (resp *types.AlbumBackVO, err error) {
	in := &resourcerpc.IdReq{
		Id: req.Id,
	}

	out, err := l.svcCtx.ResourceRpc.GetAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.AlbumBackVO{
		Id:         out.Id,
		AlbumName:  out.AlbumName,
		AlbumDesc:  out.AlbumDesc,
		AlbumCover: out.AlbumCover,
		IsDelete:   out.IsDelete,
		Status:     out.Status,
		CreatedAt:  out.CreatedAt,
		UpdatedAt:  out.UpdatedAt,
		PhotoCount: out.PhotoCount,
	}, nil
}
