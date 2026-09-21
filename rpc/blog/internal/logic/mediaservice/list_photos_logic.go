package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListPhotosLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPhotosLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPhotosLogic {
	return &ListPhotosLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPhotosLogic) ListPhotos(in *mediarpc.ListPhotosRequest) (*mediarpc.ListPhotosResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.AlbumId != 0 {
		opts = append(opts, queryx.WithCondition("album_id = ?", in.AlbumId))
	}
	if in.IsDelete != nil {
		opts = append(opts, queryx.WithCondition("is_delete = ?", *in.IsDelete))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TPhotoModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*mediarpc.Photo
	for _, v := range records {
		list = append(list, convertPhotoOut(v))
	}

	return &mediarpc.ListPhotosResponse{
		ListResult: &mediarpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
