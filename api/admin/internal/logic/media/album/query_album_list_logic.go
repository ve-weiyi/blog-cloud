package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/mediaservice"
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

func (l *QueryAlbumListLogic) QueryAlbumList(req *types.QueryAlbumListReq) (resp *types.ListResult, err error) {
	out, err := l.svcCtx.MediaService.ListAlbums(l.ctx, &mediaservice.ListAlbumsRequest{
		ListQuery: &mediaservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
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

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
