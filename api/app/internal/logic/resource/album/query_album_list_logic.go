package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/resourceservice"
)

type QueryAlbumListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取相册列表
func NewQueryAlbumListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAlbumListLogic {
	return &QueryAlbumListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryAlbumListLogic) QueryAlbumList(req *types.QueryAlbumListReq) (resp *types.PageResult, err error) {
	isDelete := int64(0)

	in := &resourceservice.ListAlbumsRequest{
		PageQuery: &resourceservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		IsDelete: &isDelete,
	}

	out, err := l.svcCtx.ResourceService.ListAlbums(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Album, 0)
	for _, v := range out.List {
		list = append(list, &types.Album{
			Id:         v.Id,
			AlbumName:  v.AlbumName,
			AlbumDesc:  v.AlbumDesc,
			AlbumCover: v.AlbumCover,
		})
	}

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
