package category

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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
	out, err := l.svcCtx.ArticleService.ListCategories(l.ctx, &articleservice.ListCategoriesRequest{
		PageQuery:    &articleservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		CategoryName: req.CategoryName,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.CategoryVO
	for _, v := range out.List {
		list = append(list, &types.CategoryVO{
			Id:           v.Id,
			CategoryName: v.CategoryName,
			ArticleCount: v.ArticleCount,
			CreatedAt:    v.CreatedAt,
			UpdatedAt:    v.UpdatedAt,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
