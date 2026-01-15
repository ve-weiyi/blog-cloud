package album

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAlbumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新相册
func NewUpdateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAlbumLogic {
	return &UpdateAlbumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAlbumLogic) UpdateAlbum(req *types.NewAlbumReq) (resp *types.AlbumBackVO, err error) {
	in := &resourcerpc.UpdateAlbumReq{
		Id:         req.Id,
		AlbumName:  req.AlbumName,
		AlbumDesc:  req.AlbumDesc,
		AlbumCover: req.AlbumCover,
		IsDelete:   req.IsDelete,
		Status:     req.Status,
	}

	out, err := l.svcCtx.ResourceRpc.UpdateAlbum(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertAlbumTypes(out.Album), nil
}
