package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type QueryArchivedArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取归档文章列表
func NewQueryArchivedArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryArchivedArticleListLogic {
	return &QueryArchivedArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryArchivedArticleListLogic) QueryArchivedArticleList(req *types.QueryArchivedArticleListReq) (resp *types.PageResult, err error) {
	isDelete := int64(0)
	status := int64(1)

	in := &articleservice.ListArticlesRequest{
		PageQuery: &articleservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    []string{"created_at desc"},
		},
		IsDelete: &isDelete,
		Status:   &status,
	}

	out, err := l.svcCtx.ArticleService.ListArticles(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.ArticleHome, 0)
	for _, v := range out.List {
		m := convertArticleHomeTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}
	return
}
