package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
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

func (l *ListPhotosLogic) ListPhotos(in *resourcerpc.ListPhotosRequest) (*resourcerpc.ListPhotosResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
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

	var list []*resourcerpc.Photo
	for _, v := range records {
		list = append(list, convertPhotoOut(v))
	}

	return &resourcerpc.ListPhotosResponse{
		PageResult: &resourcerpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
