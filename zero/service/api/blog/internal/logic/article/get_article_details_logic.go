package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"
)

type GetArticleDetailsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文章详情
func NewGetArticleDetailsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleDetailsLogic {
	return &GetArticleDetailsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleDetailsLogic) GetArticleDetails(req *types.IdReq) (resp *types.ArticleDeatils, err error) {
	in := &articlerpc.IdReq{
		Id: req.Id,
	}

	out, err := l.svcCtx.ArticleRpc.GetArticle(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询上一篇文章
	recommend, err := l.svcCtx.ArticleRpc.GetArticleRecommend(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.ArticleDeatils{
		ArticleHome:          types.ArticleHome{},
		LastArticle:          nil,
		NextArticle:          nil,
		RecommendArticleList: nil,
		NewestArticleList:    nil,
	}

	resp.ArticleHome = *ConvertArticleHomeTypes(out)

	resp.LastArticle = ConvertArticlePreviewTypes(recommend.Last)

	resp.NextArticle = ConvertArticlePreviewTypes(recommend.Next)

	for _, v := range recommend.Recommend {
		resp.RecommendArticleList = append(resp.RecommendArticleList, ConvertArticlePreviewTypes(v))
	}

	for _, v := range recommend.Newest {
		resp.NewestArticleList = append(resp.NewestArticleList, ConvertArticlePreviewTypes(v))
	}

	return
}
