package job

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/cachekey"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/domain/knowledge"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	// articleChunkInterval 增量重建周期。
	articleChunkInterval = 5 * time.Minute
	// articleChunkFirstDelay 启动后的首轮延迟，避开启动流量。
	articleChunkFirstDelay = time.Minute
	// articleChunkBatchSize 每批处理的文章数。
	articleChunkBatchSize = 100
)

// ArticleChunkIndexer 把公开文章切成分块，供 AI 助手检索。
//
// 增量以 (updated_at, id) 为游标推进，位点存 Redis：文章转为非公开或软删除时
// 同一条流程删除分块，孤儿分块每轮兜底清理。
type ArticleChunkIndexer struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewArticleChunkIndexer(db *gorm.DB, rdb *redis.Client) *ArticleChunkIndexer {
	return &ArticleChunkIndexer{db: db, rdb: rdb}
}

// Task 返回可注册到 Runner 的后台任务。
func (j *ArticleChunkIndexer) Task() Task {
	periodic := NewPeriodic("article-chunk", articleChunkInterval, articleChunkFirstDelay, func(ctx context.Context) {
		if err := j.RunOnce(ctx); err != nil {
			logx.WithContext(ctx).Errorf("文章分块重建失败: %v", err)
		}
	})
	return Task{Name: "article-chunk", Start: periodic.Start, Stop: periodic.Stop}
}

// RunOnce 处理一批增量文章并清理孤儿分块。
func (j *ArticleChunkIndexer) RunOnce(ctx context.Context) error {
	release, ok, err := AcquireLock(ctx, j.db, "blog:article-chunk")
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	defer release()

	cursor, err := j.loadCursor(ctx)
	if err != nil {
		return err
	}
	if cursor.UpdatedAt.IsZero() {
		// MySQL datetime 不接受公元 1 年，用 Unix 纪元作为空位点
		cursor.UpdatedAt = time.Unix(0, 0)
	}

	articleModel := model.NewTArticleModel(j.db)
	for {
		list, _, err := articleModel.FindListAndTotal(ctx, 1, articleChunkBatchSize, "updated_at asc, id asc",
			"updated_at > ? OR (updated_at = ? AND id > ?)", cursor.UpdatedAt, cursor.UpdatedAt, cursor.Id)
		if err != nil {
			return fmt.Errorf("job: list updated articles: %w", err)
		}
		if len(list) == 0 {
			break
		}

		for _, article := range list {
			if err := j.rebuild(ctx, article); err != nil {
				return fmt.Errorf("job: rebuild article %d: %w", article.Id, err)
			}
			cursor = chunkCursor{UpdatedAt: article.UpdatedAt, Id: article.Id}
			if err := j.saveCursor(ctx, cursor); err != nil {
				return err
			}
		}
		if len(list) < articleChunkBatchSize {
			break
		}
	}

	return j.cleanOrphans(ctx)
}

// rebuild 重建单篇文章的分块；非公开或已删除的文章只删块。
func (j *ArticleChunkIndexer) rebuild(ctx context.Context, article *model.TArticle) error {
	return j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		chunkModel := model.NewTArticleChunkModel(tx)

		if article.IsDelete != 0 || article.Status != enums.ArticleStatusPublic {
			_, err := chunkModel.DeleteBatch(ctx, "article_id = ?", article.Id)
			return err
		}

		chunks := knowledge.Split(article.ArticleContent)
		next := make([]*model.TArticleChunk, 0, len(chunks))
		for _, chunk := range chunks {
			next = append(next, &model.TArticleChunk{
				ArticleId:    article.Id,
				ArticleTitle: article.ArticleTitle,
				Ordinal:      int64(chunk.Ordinal),
				Heading:      chunk.Heading,
				Content:      chunk.Content,
				Summary:      chunk.Summary,
				ContentHash:  chunk.Hash,
				UpdatedAt:    time.Now(),
			})
		}

		existing, err := chunkModel.FindAll(ctx, "article_id = ?", article.Id)
		if err != nil {
			return err
		}
		if sameChunks(existing, next) {
			// 内容没有实质变化：跳过写操作
			return nil
		}

		if _, err := chunkModel.DeleteBatch(ctx, "article_id = ?", article.Id); err != nil {
			return err
		}
		if len(next) == 0 {
			return nil
		}
		_, err = chunkModel.InsertBatch(ctx, next...)
		return err
	})
}

// sameChunks 比较分块内容是否需要重建；标题变化也算变化。
func sameChunks(existing, next []*model.TArticleChunk) bool {
	if len(existing) != len(next) {
		return false
	}
	known := make(map[int64]string, len(existing))
	for _, chunk := range existing {
		known[chunk.Ordinal] = chunk.ContentHash + "|" + chunk.ArticleTitle
	}
	for _, chunk := range next {
		if known[chunk.Ordinal] != chunk.ContentHash+"|"+chunk.ArticleTitle {
			return false
		}
	}
	return true
}

// cleanOrphans 删除文章已不存在或已不可见的分块，兜底处理重建窗口与历史残留。
func (j *ArticleChunkIndexer) cleanOrphans(ctx context.Context) error {
	const query = `
DELETE c FROM t_article_chunk c
LEFT JOIN t_article a ON a.id = c.article_id
WHERE a.id IS NULL
   OR a.is_delete = 1
   OR a.status <> ?`

	if err := j.db.WithContext(ctx).Exec(query, enums.ArticleStatusPublic).Error; err != nil {
		return fmt.Errorf("job: clean orphan chunks: %w", err)
	}
	return nil
}

type chunkCursor struct {
	UpdatedAt time.Time
	Id        int64
}

func (c chunkCursor) encode() string {
	return fmt.Sprintf("%d|%d", c.UpdatedAt.Unix(), c.Id)
}

func decodeChunkCursor(raw string) chunkCursor {
	parts := strings.SplitN(raw, "|", 2)
	if len(parts) != 2 {
		return chunkCursor{}
	}
	sec, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return chunkCursor{}
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return chunkCursor{}
	}
	return chunkCursor{UpdatedAt: time.Unix(sec, 0), Id: id}
}

func (j *ArticleChunkIndexer) loadCursor(ctx context.Context) (chunkCursor, error) {
	if j.rdb == nil {
		return chunkCursor{}, fmt.Errorf("job: redis is required for article chunk cursor")
	}
	raw, err := j.rdb.Get(ctx, cachekey.ArticleChunkWatermarkKey).Result()
	if err != nil {
		if err == redis.Nil {
			return chunkCursor{}, nil
		}
		return chunkCursor{}, fmt.Errorf("job: load article chunk cursor: %w", err)
	}
	return decodeChunkCursor(raw), nil
}

func (j *ArticleChunkIndexer) saveCursor(ctx context.Context, cursor chunkCursor) error {
	if err := j.rdb.Set(ctx, cachekey.ArticleChunkWatermarkKey, cursor.encode(), 0).Err(); err != nil {
		return fmt.Errorf("job: save article chunk cursor: %w", err)
	}
	return nil
}
