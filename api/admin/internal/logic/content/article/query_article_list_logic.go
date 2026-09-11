package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
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
	out, err := l.svcCtx.ArticleService.ListArticles(l.ctx, &articleservice.ListArticlesRequest{
		PageQuery:    &articleservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		ArticleTitle: req.ArticleTitle,
		ArticleType:  req.ArticleType,
		CategoryName: req.CategoryName,
		TagName:      req.TagName,
		IsTop:        req.IsTop,
		IsDelete:     req.IsDelete,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.ArticleVO
	for _, v := range out.List {
		var tagNames []string
		for _, t := range v.Tags {
			tagNames = append(tagNames, t.TagName)
		}
		categoryName := ""
		if v.Category != nil {
			categoryName = v.Category.CategoryName
		}
		list = append(list, &types.ArticleVO{
			Id:             v.Id,
			ArticleCover:   v.ArticleCover,
			ArticleTitle:   v.ArticleTitle,
			ArticleContent: v.ArticleContent,
			ArticleType:    v.ArticleType,
			OriginalUrl:    v.OriginalUrl,
			IsTop:          v.IsTop,
			IsDelete:       v.IsDelete,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			CategoryName:   categoryName,
			TagNameList:    tagNames,
			LikeCount:      v.LikeCount,
			ViewsCount:     v.ViewCount,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
