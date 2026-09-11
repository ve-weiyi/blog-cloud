package analyticsservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/constants/cachekey"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/analyticsrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type GetArticleStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetArticleStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetArticleStatsLogic {
	return &GetArticleStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文章数据分析（从业务表查询）
func (l *GetArticleStatsLogic) GetArticleStats(in *analyticsrpc.GetArticleStatsRequest) (*analyticsrpc.GetArticleStatsResponse, error) {
	categoryCount, err := l.svcCtx.TCategoryModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	tagCount, err := l.svcCtx.TTagModel.FindCount(l.ctx, "")
	if err != nil {
		return nil, err
	}

	categoryStats, err := l.getCategoryStats()
	if err != nil {
		return nil, err
	}

	tagStats, err := l.getTagStats()
	if err != nil {
		return nil, err
	}

	articleRanks, err := l.getArticleRanks(10)
	if err != nil {
		return nil, err
	}

	dailyStats, err := l.getDailyStats()
	if err != nil {
		return nil, err
	}

	return &analyticsrpc.GetArticleStatsResponse{
		CategoryCount:   categoryCount,
		TagCount:        tagCount,
		CategoryStats:   categoryStats,
		TagStats:        tagStats,
		ArticleRanks:    articleRanks,
		DailyStatistics: dailyStats,
	}, nil
}

func (l *GetArticleStatsLogic) getCategoryStats() ([]*analyticsrpc.CategoryStat, error) {
	var results []struct {
		CategoryID   int64 `gorm:"column:category_id"`
		ArticleCount int64 `gorm:"column:article_count"`
	}
	err := l.svcCtx.GormDB.Model(&model.TArticle{}).
		Select("category_id, COUNT(*) as article_count").
		Group("category_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64)
	for _, r := range results {
		countMap[r.CategoryID] = r.ArticleCount
	}

	categories, err := l.svcCtx.TCategoryModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	list := make([]*analyticsrpc.CategoryStat, 0, len(categories))
	for _, c := range categories {
		list = append(list, &analyticsrpc.CategoryStat{
			Id:           c.Id,
			CategoryName: c.CategoryName,
			ArticleCount: countMap[c.Id],
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getTagStats() ([]*analyticsrpc.TagStat, error) {
	var results []struct {
		TagID        int64 `gorm:"column:tag_id"`
		ArticleCount int64 `gorm:"column:article_count"`
	}
	err := l.svcCtx.GormDB.Model(&model.TArticleTag{}).
		Select("tag_id, COUNT(*) as article_count").
		Group("tag_id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64)
	for _, r := range results {
		countMap[r.TagID] = r.ArticleCount
	}

	tags, err := l.svcCtx.TTagModel.FindALL(l.ctx, "")
	if err != nil {
		return nil, err
	}

	list := make([]*analyticsrpc.TagStat, 0, len(tags))
	for _, t := range tags {
		list = append(list, &analyticsrpc.TagStat{
			Id:           t.Id,
			TagName:      t.TagName,
			ArticleCount: countMap[t.Id],
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getArticleRanks(limit int64) ([]*analyticsrpc.ArticleRank, error) {
	key := cachekey.ArticleViewCountKey
	ids, err := l.svcCtx.Redis.ZRevRange(l.ctx, key, 0, limit).Result()
	if err != nil {
		return nil, err
	}

	list := make([]*analyticsrpc.ArticleRank, 0, len(ids))
	for _, id := range ids {
		articleId := cast.ToInt64(id)
		article, err := l.svcCtx.TArticleModel.FindById(l.ctx, articleId)
		if err != nil {
			continue
		}
		viewCount, _ := l.svcCtx.Redis.ZScore(l.ctx, key, id).Result()
		list = append(list, &analyticsrpc.ArticleRank{
			Id:           article.Id,
			ArticleTitle: article.ArticleTitle,
			ViewCount:    int64(viewCount),
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getDailyStats() ([]*analyticsrpc.ArticleDailyStatistic, error) {
	var results []struct {
		Date         string `gorm:"column:date"`
		ArticleCount int64  `gorm:"column:article_count"`
	}
	err := l.svcCtx.GormDB.Raw("SELECT DATE(created_at) AS date, COUNT(*) as article_count FROM t_article GROUP BY date ORDER BY date DESC").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	list := make([]*analyticsrpc.ArticleDailyStatistic, 0, len(results))
	for _, r := range results {
		list = append(list, &analyticsrpc.ArticleDailyStatistic{
			Date:  r.Date,
			Count: r.ArticleCount,
		})
	}
	return list, nil
}
