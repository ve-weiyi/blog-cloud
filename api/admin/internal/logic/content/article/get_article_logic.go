package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type GetArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文章详情
func NewGetArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleLogic {
	return &GetArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleLogic) GetArticle(req *types.GetArticleReq) (resp *types.ArticleVO, err error) {
	out, err := l.svcCtx.ArticleService.GetArticle(l.ctx, &articleservice.GetArticleRequest{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	article := out.Article
	var tagNames []string
	for _, t := range article.Tags {
		tagNames = append(tagNames, t.TagName)
	}
	categoryName := ""
	if article.Category != nil {
		categoryName = article.Category.CategoryName
	}
	return &types.ArticleVO{
		Id:             article.Id,
		ArticleCover:   article.ArticleCover,
		ArticleTitle:   article.ArticleTitle,
		ArticleContent: article.ArticleContent,
		ArticleType:    article.ArticleType,
		OriginalUrl:    article.OriginalUrl,
		IsTop:          article.IsTop,
		IsDelete:       article.IsDelete,
		Status:         article.Status,
		CreatedAt:      article.CreatedAt,
		UpdatedAt:      article.UpdatedAt,
		CategoryName:   categoryName,
		TagNameList:    tagNames,
		LikeCount:      article.LikeCount,
		ViewsCount:     article.ViewCount,
	}, nil
}
