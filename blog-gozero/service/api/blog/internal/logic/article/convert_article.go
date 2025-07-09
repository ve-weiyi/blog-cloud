package article

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
)

func ConvertArticleHomeTypes(in *articlerpc.ArticleDetailsResp) (out *types.ArticleHome) {
	if in == nil {
		return nil
	}

	var tags []string
	for _, tag := range in.TagList {
		tags = append(tags, tag.TagName)
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
		CategoryName:   in.Category.CategoryName,
		TagNameList:    tags,
		LikeCount:      in.LikeCount,
		ViewsCount:     in.ViewCount,
	}

	return
}

func ConvertArticlePreviewTypes(in *articlerpc.ArticlePreview) (out *types.ArticlePreview) {
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
