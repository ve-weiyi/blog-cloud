package article

import (
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"
)

func ConvertArticleHomeTypes(in *articlerpc.ArticleDetails) (out *types.ArticleHome) {
	if in == nil {
		return nil
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
		CategoryName:   "",
		TagNameList:    make([]string, 0),
		LikeCount:      in.LikeCount,
		ViewsCount:     0,
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
		CreatedAt:    in.CreatedAt,
	}
	return
}
