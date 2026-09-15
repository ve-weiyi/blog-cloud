package overview

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/statsservice"
)

type GetArticleStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文章分析数据
func NewGetArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleStatsLogic {
	return &GetArticleStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetArticleStatsLogic) GetArticleStats(req *types.EmptyReq) (resp *types.GetArticleStatsResp, err error) {
	out, err := l.svcCtx.StatsService.GetArticleStats(l.ctx, &statsservice.GetArticleStatsRequest{})
	if err != nil {
		return nil, err
	}

	resp = &types.GetArticleStatsResp{
		CategoryList:      make([]*types.CategoryOverviewVO, 0, len(out.CategoryStats)),
		TagList:           make([]*types.TagOverviewVO, 0, len(out.TagStats)),
		ArticleViewRanks:  make([]*types.ArticleViewVO, 0, len(out.ArticleRanks)),
		ArticleStatistics: make([]*types.ArticleStatisticsVO, 0, len(out.DailyStatistics)),
	}

	for _, c := range out.CategoryStats {
		resp.CategoryList = append(resp.CategoryList, &types.CategoryOverviewVO{
			Id:           c.Id,
			CategoryName: c.CategoryName,
			ArticleCount: c.ArticleCount,
		})
	}
	for _, t := range out.TagStats {
		resp.TagList = append(resp.TagList, &types.TagOverviewVO{
			Id:           t.Id,
			TagName:      t.TagName,
			ArticleCount: t.ArticleCount,
		})
	}
	for _, r := range out.ArticleRanks {
		resp.ArticleViewRanks = append(resp.ArticleViewRanks, &types.ArticleViewVO{
			Id:           r.Id,
			ArticleTitle: r.ArticleTitle,
			ViewCount:    r.ViewCount,
		})
	}
	for _, s := range out.DailyStatistics {
		resp.ArticleStatistics = append(resp.ArticleStatistics, &types.ArticleStatisticsVO{
			Date:  s.Date,
			Count: s.Count,
		})
	}

	return
}
