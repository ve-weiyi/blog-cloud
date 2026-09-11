package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListAlbumsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAlbumsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAlbumsLogic {
	return &ListAlbumsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAlbumsLogic) ListAlbums(in *resourcerpc.ListAlbumsRequest) (*resourcerpc.ListAlbumsResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if in.IsDelete != nil {
		opts = append(opts, queryx.WithCondition("is_delete = ?", *in.IsDelete))
	}
	if in.AlbumName != nil {
		opts = append(opts, queryx.WithCondition("album_name like ?", "%"+*in.AlbumName+"%"))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TAlbumModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	cm := l.findPhotoCountGroupAlbum(records)

	var list []*resourcerpc.Album
	for _, v := range records {
		list = append(list, convertAlbumOut(v, cm))
	}

	return &resourcerpc.ListAlbumsResponse{
		PageResult: &resourcerpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}

func (l *ListAlbumsLogic) findPhotoCountGroupAlbum(list []*model.TAlbum) map[int64]int {
	var ids []int64
	for _, v := range list {
		ids = append(ids, v.Id)
	}
	if len(ids) == 0 {
		return nil
	}

	var results []struct {
		AlbumID    int64 `gorm:"column:album_id"`
		PhotoCount int   `gorm:"column:photo_count"`
	}
	err := l.svcCtx.GormDB.Model(&model.TPhoto{}).
		Select("album_id, COUNT(*) as photo_count").
		Where("album_id IN ?", ids).
		Group("album_id").
		Scan(&results).Error
	if err != nil {
		return nil
	}

	cm := make(map[int64]int)
	for _, r := range results {
		cm[r.AlbumID] = r.PhotoCount
	}
	return cm
}
