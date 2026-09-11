package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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
	out, err := l.svcCtx.ResourceService.ListAlbums(l.ctx, &resourceservice.ListAlbumsRequest{
		PageQuery: &resourceservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		AlbumName: req.AlbumName,
		IsDelete:  req.IsDelete,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.AlbumVO
	for _, v := range out.List {
		list = append(list, &types.AlbumVO{
			Id:         v.Id,
			AlbumName:  v.AlbumName,
			AlbumDesc:  v.AlbumDesc,
			AlbumCover: v.AlbumCover,
			IsDelete:   v.IsDelete,
			Status:     v.Status,
			PhotoCount: v.PhotoCount,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
