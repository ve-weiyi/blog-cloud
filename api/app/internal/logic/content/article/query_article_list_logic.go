package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type QueryArticleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文章列表
func NewQueryArticleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryArticleListLogic {
	return &QueryArticleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryArticleListLogic) QueryArticleList(req *types.QueryArticleListReq) (resp *types.PageResult, err error) {
	isDelete := int64(0)
	status := int64(1)

	in := &articleservice.ListArticlesRequest{
		PageQuery: &articleservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		ArticleTitle: req.ArticleTitle,
		ArticleType:  nil,
		CategoryName: req.CategoryName,
		TagName:      req.TagName,
		IsTop:        nil,
		IsDelete:     &isDelete,
		Status:       &status,
		Ids:          nil,
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
