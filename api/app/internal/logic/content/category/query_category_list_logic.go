package category

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type QueryCategoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取分类列表
func NewQueryCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCategoryListLogic {
	return &QueryCategoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCategoryListLogic) QueryCategoryList(req *types.QueryCategoryListReq) (resp *types.PageResult, err error) {
	in := &articleservice.ListCategoriesRequest{
		PageQuery: &articleservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		CategoryName: req.CategoryName,
	}

	out, err := l.svcCtx.ArticleService.ListCategories(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Category, 0)
	for _, v := range out.List {
		list = append(list, &types.Category{
			Id:           v.Id,
			CategoryName: v.CategoryName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
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
