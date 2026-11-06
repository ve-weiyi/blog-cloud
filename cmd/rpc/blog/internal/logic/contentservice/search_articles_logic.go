package contentservicelogic

import (
	"context"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
)

const (
	// defaultSearchChunks 未指定时的片段数量。
	defaultSearchChunks = 5
	// maxSearchChunks 单次检索返回的片段上限。
	maxSearchChunks = 20
	// maxKeywordRunes 关键词的最大字符数，避免把整段提问当检索词。
	maxKeywordRunes = 100
	// maxChunksPerArticle 每篇文章最多保留的片段数，防止一篇长文占满上下文。
	maxChunksPerArticle = 2
	// searchFetchMultiplier 为去重预留余量：按目标数量的倍数取原始命中。
	searchFetchMultiplier = 4
)

type SearchArticlesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchArticlesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchArticlesLogic {
	return &SearchArticlesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检索公开文章分块（供 AI 助手使用）
func (l *SearchArticlesLogic) SearchArticles(in *contentrpc.SearchArticlesRequest) (*contentrpc.SearchArticlesResponse, error) {
	keyword := truncateRunes(strings.TrimSpace(in.Keyword), maxKeywordRunes)
	if keyword == "" {
		return &contentrpc.SearchArticlesResponse{}, nil
	}

	limit := int(in.Limit)
	if limit <= 0 {
		limit = defaultSearchChunks
	}
	if limit > maxSearchChunks {
		limit = maxSearchChunks
	}

	hits, err := l.svcCtx.TArticleChunkModel.Search(l.ctx, keyword, enums.ArticleStatusPublic, limit*searchFetchMultiplier)
	if err != nil {
		return nil, err
	}

	list := make([]*contentrpc.ArticleChunk, 0, limit)
	for _, hit := range dedupeChunkHits(hits, maxChunksPerArticle, limit) {
		list = append(list, &contentrpc.ArticleChunk{
			ArticleId:    hit.ArticleId,
			ArticleTitle: hit.ArticleTitle,
			Heading:      hit.Heading,
			Content:      hit.Content,
			Score:        hit.Score,
			Summary:      hit.Summary,
		})
	}

	return &contentrpc.SearchArticlesResponse{List: list}, nil
}

// dedupeChunkHits 按得分排序，并限制同一篇文章的片段数与总片段数。
func dedupeChunkHits(hits []*model.ArticleChunkHit, perArticle, limit int) []*model.ArticleChunkHit {
	if len(hits) == 0 || limit <= 0 {
		return nil
	}

	sorted := make([]*model.ArticleChunkHit, len(hits))
	copy(sorted, hits)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})

	counts := make(map[int64]int, limit)
	selected := make([]*model.ArticleChunkHit, 0, limit)
	for _, hit := range sorted {
		if counts[hit.ArticleId] >= perArticle {
			continue
		}
		counts[hit.ArticleId]++
		selected = append(selected, hit)
		if len(selected) >= limit {
			break
		}
	}
	return selected
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}
