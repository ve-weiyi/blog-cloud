package article

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
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

func (l *GetArticleLogic) GetArticle(req *types.GetArticleReq) (resp *types.ArticleDetails, err error) {
	_, err = l.svcCtx.ArticleService.IncrementArticleView(l.ctx, &articleservice.IncrementArticleViewRequest{
		Id: req.ArticleId,
	})
	if err != nil {
		return nil, err
	}

	articleOut, err := l.svcCtx.ArticleService.GetArticle(l.ctx, &articleservice.GetArticleRequest{
		Id: req.ArticleId,
	})
	if err != nil {
		return nil, err
	}

	relationOut, err := l.svcCtx.ArticleService.GetArticleRelation(l.ctx, &articleservice.GetArticleRelationRequest{
		Id: req.ArticleId,
	})
	if err != nil {
		return nil, err
	}

	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, []string{articleOut.Article.UserId})
	if err != nil {
		return nil, err
	}

	resp = &types.ArticleDetails{
		ArticleHome:          *convertArticleHomeTypes(articleOut.Article),
		Author:               usm[articleOut.Article.UserId],
		LastArticle:          convertArticlePreviewTypes(relationOut.Last),
		NextArticle:          convertArticlePreviewTypes(relationOut.Next),
		RecommendArticleList: make([]*types.ArticlePreview, 0),
		NewestArticleList:    make([]*types.ArticlePreview, 0),
	}

	for _, v := range relationOut.Recommends {
		resp.RecommendArticleList = append(resp.RecommendArticleList, convertArticlePreviewTypes(v))
	}

	for _, v := range relationOut.Newests {
		resp.NewestArticleList = append(resp.NewestArticleList, convertArticlePreviewTypes(v))
	}

	return
}
