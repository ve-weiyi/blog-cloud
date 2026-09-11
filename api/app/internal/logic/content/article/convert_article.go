package article

import (
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

func convertArticleHomeTypes(in *articleservice.Article) (out *types.ArticleHome) {
	if in == nil {
		return nil
	}

	var tags []string
	for _, tag := range in.Tags {
		tags = append(tags, tag.TagName)
	}

	var categoryName string
	if in.Category != nil {
		categoryName = in.Category.CategoryName
	}

	out = &types.ArticleHome{
		Id:             in.Id,
		ArticleCover:   in.ArticleCover,
		ArticleTitle:   in.ArticleTitle,
		ArticleContent: in.ArticleContent,
		ArticleType:    in.ArticleType,
		OriginalUrl:    in.OriginalUrl,
		IsTop:          in.IsTop,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		CategoryName:   categoryName,
		TagNameList:    tags,
		LikeCount:      in.LikeCount,
		ViewsCount:     in.ViewCount,
	}

	return
}

func convertArticlePreviewTypes(in *articleservice.ArticlePreview) (out *types.ArticlePreview) {
	if in == nil {
		return nil
	}

	out = &types.ArticlePreview{
		Id:           in.Id,
		ArticleCover: in.ArticleCover,
		ArticleTitle: in.ArticleTitle,
		LikeCount:    in.LikeCount,
		ViewsCount:   in.ViewCount,
		CreatedAt:    in.CreatedAt,
	}
	return
}
