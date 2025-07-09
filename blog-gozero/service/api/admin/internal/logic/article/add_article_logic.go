package article

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

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

func (l *AddArticleLogic) AddArticle(req *types.ArticleNewReq) (resp *types.ArticleBackVO, err error) {
	in := &articlerpc.ArticleNewReq{
		Id:             req.Id,
		UserId:         cast.ToString(l.ctx.Value(restx.HeaderUid)),
		ArticleCover:   req.ArticleCover,
		ArticleTitle:   req.ArticleTitle,
		ArticleContent: req.ArticleContent,
		ArticleType:    req.ArticleType,
		OriginalUrl:    req.OriginalUrl,
		IsTop:          req.IsTop,
		Status:         req.Status,
		CategoryName:   req.CategoryName,
		TagNameList:    req.TagNameList,
	}

	out, err := l.svcCtx.ArticleRpc.AddArticle(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return &types.ArticleBackVO{
		Id: out.Id,
	}, nil
}

func ConvertArticleTypes(in *articlerpc.ArticleDetailsResp) (out *types.ArticleBackVO) {
	out = &types.ArticleBackVO{
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
