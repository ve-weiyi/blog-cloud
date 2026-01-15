package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AnalysisArticleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAnalysisArticleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AnalysisArticleLogic {
	return &AnalysisArticleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询文章数量
func (l *AnalysisArticleLogic) AnalysisArticle(in *articlerpc.AnalysisArticleReq) (*articlerpc.AnalysisArticleResp, error) {
	ac, err := l.svcCtx.TArticleModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	cc, err := l.svcCtx.TCategoryModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	tc, err := l.svcCtx.TTagModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	helper := NewArticleHelperLogic(l.ctx, l.svcCtx)

	cl, err := l.svcCtx.TCategoryModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	cds, err := helper.convertCategory(cl)
	if err != nil {
		return nil, err
	}

	tl, err := l.svcCtx.TTagModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	tds, err := helper.convertTag(tl)
	if err != nil {
		return nil, err
	}

	tops, err := helper.GetViewTopArticleList(10)
	if err != nil {
		return nil, err
	}

	var ars []*articlerpc.ArticlePreview
	for _, article := range tops {
		m := helper.convertArticlePreviewOut(article)
		ars = append(ars, m)
	}

	daily, err := helper.GetArticleDailyStatistics()
	if err != nil {
		return nil, err
	}

	var ads []*articlerpc.ArticleDailyStatistics
	for k, v := range daily {
		m := &articlerpc.ArticleDailyStatistics{
			Date:  k,
			Count: v,
		}
		ads = append(ads, m)
	}

	out := &articlerpc.AnalysisArticleResp{
		ArticleCount:           ac,
		CategoryCount:          cc,
		TagCount:               tc,
		CategoryList:           cds,
		TagList:                tds,
		ArticleRankList:        ars,
		ArticleDailyStatistics: ads,
	}

	return out, nil
}
