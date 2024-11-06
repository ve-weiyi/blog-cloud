package article

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddArticleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加文章
func NewAddArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddArticleLogic {
	return &AddArticleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddArticleLogic) AddArticle(req *types.ArticleNewReq) (resp *types.ArticleBackDTO, err error) {
	in := ConvertArticlePb(req)
	in.UserId = cast.ToString(l.ctx.Value("uid"))
	out, err := l.svcCtx.ArticleRpc.AddArticle(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertArticleTypes(out)
	return resp, nil
}

func ConvertArticlePb(in *types.ArticleNewReq) (out *articlerpc.ArticleNewReq) {
	out = &articlerpc.ArticleNewReq{
		Id:             in.Id,
		UserId:         "",
		ArticleCover:   in.ArticleCover,
		ArticleTitle:   in.ArticleTitle,
		ArticleContent: in.ArticleContent,
		ArticleType:    in.ArticleType,
		OriginalUrl:    in.OriginalUrl,
		Status:         in.Status,
		CategoryName:   in.CategoryName,
		TagNameList:    in.TagNameList,
	}

	return
}

func ConvertArticleTypes(in *articlerpc.ArticleDetails) (out *types.ArticleBackDTO) {
	out = &types.ArticleBackDTO{
		Id:             in.Id,
		ArticleCover:   in.ArticleCover,
		ArticleTitle:   in.ArticleTitle,
		ArticleContent: in.ArticleContent,
		ArticleType:    in.ArticleType,
		OriginalUrl:    in.OriginalUrl,
		IsTop:          in.IsTop,
		IsDelete:       in.IsDelete,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
		CategoryName:   "",
		TagNameList:    make([]string, 0),
		LikeCount:      in.LikeCount,
		ViewsCount:     in.ViewCount,
	}

	if in.Category != nil {
		out.CategoryName = in.Category.CategoryName
	}

	if in.TagList != nil {
		for _, tag := range in.TagList {
			out.TagNameList = append(out.TagNameList, tag.TagName)
		}
	}

	return
}
