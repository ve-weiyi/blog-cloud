package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCategoriesLogic) ListCategories(in *articlerpc.ListCategoriesRequest) (*articlerpc.ListCategoriesResponse, error) {
	helper := NewArticleHelper(l.ctx, l.svcCtx)

	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if in.CategoryName != nil {
		opts = append(opts, queryx.WithCondition("category_name like ?", "%"+*in.CategoryName+"%"))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TCategoryModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	acm, _ := helper.findArticleCountGroupCategory(records)

	var list []*articlerpc.Category
	for _, entity := range records {
		m := &articlerpc.Category{
			Id:           entity.Id,
			CategoryName: entity.CategoryName,
			CreatedAt:    entity.CreatedAt.UnixMilli(),
			UpdatedAt:    entity.UpdatedAt.UnixMilli(),
		}
		if acm != nil {
			m.ArticleCount = acm[entity.Id]
		}
		list = append(list, m)
	}

	return &articlerpc.ListCategoriesResponse{
		PageResult: &articlerpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
