package statsservicelogic

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/statsrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/model"
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
func (l *GetArticleStatsLogic) GetArticleStats(in *statsrpc.GetArticleStatsRequest) (*statsrpc.GetArticleStatsResponse, error) {
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

	return &statsrpc.GetArticleStatsResponse{
		CategoryCount:   categoryCount,
		TagCount:        tagCount,
		CategoryStats:   categoryStats,
		TagStats:        tagStats,
		ArticleRanks:    articleRanks,
		DailyStatistics: dailyStats,
	}, nil
}

func (l *GetArticleStatsLogic) getCategoryStats() ([]*statsrpc.CategoryStat, error) {
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

	categories, err := l.svcCtx.TCategoryModel.FindAll(l.ctx, "")
	if err != nil {
		return nil, err
	}

	list := make([]*statsrpc.CategoryStat, 0, len(categories))
	for _, c := range categories {
		list = append(list, &statsrpc.CategoryStat{
			Id:           c.Id,
			CategoryName: c.CategoryName,
			ArticleCount: countMap[c.Id],
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getTagStats() ([]*statsrpc.TagStat, error) {
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

	tags, err := l.svcCtx.TTagModel.FindAll(l.ctx, "")
	if err != nil {
		return nil, err
	}

	list := make([]*statsrpc.TagStat, 0, len(tags))
	for _, t := range tags {
		list = append(list, &statsrpc.TagStat{
			Id:           t.Id,
			TagName:      t.TagName,
			ArticleCount: countMap[t.Id],
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getArticleRanks(limit int64) ([]*statsrpc.ArticleRank, error) {
	key := cachekey.ArticleViewCountKey
	items, err := l.svcCtx.Redis.ZRevRangeWithScores(l.ctx, key, 0, limit).Result()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []*statsrpc.ArticleRank{}, nil
	}

	ids := make([]int64, 0, len(items))
	viewCounts := make(map[int64]int64, len(items))
	for _, item := range items {
		id := cast.ToInt64(item.Member)
		ids = append(ids, id)
		viewCounts[id] = int64(item.Score)
	}

	articles, err := l.svcCtx.TArticleModel.FindAll(l.ctx, "id IN ?", ids)
	if err != nil {
		return nil, err
	}
	articleMap := make(map[int64]*model.TArticle, len(articles))
	for _, article := range articles {
		articleMap[article.Id] = article
	}

	list := make([]*statsrpc.ArticleRank, 0, len(items))
	for _, id := range ids {
		article, ok := articleMap[id]
		if !ok {
			continue
		}
		list = append(list, &statsrpc.ArticleRank{
			Id:           article.Id,
			ArticleTitle: article.ArticleTitle,
			ViewCount:    viewCounts[id],
		})
	}
	return list, nil
}

func (l *GetArticleStatsLogic) getDailyStats() ([]*statsrpc.ArticleDailyStatistic, error) {
	var results []struct {
		Date         string `gorm:"column:date"`
		ArticleCount int64  `gorm:"column:article_count"`
	}
	err := l.svcCtx.GormDB.Raw("SELECT DATE(created_at) AS date, COUNT(*) as article_count FROM t_article GROUP BY date ORDER BY date DESC").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	list := make([]*statsrpc.ArticleDailyStatistic, 0, len(results))
	for _, r := range results {
		list = append(list, &statsrpc.ArticleDailyStatistic{
			Date:  r.Date,
			Count: r.ArticleCount,
		})
	}
	return list, nil
}
