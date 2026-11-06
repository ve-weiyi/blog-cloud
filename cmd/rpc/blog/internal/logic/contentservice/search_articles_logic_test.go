package contentservicelogic

import (
	"testing"

	"github.com/ve-weiyi/blog-cloud/model"
)

// 结果按得分排序，同一篇文章最多保留两段，总数受 limit 约束。
func TestDedupeChunkHits(t *testing.T) {
	hits := []*model.ArticleChunkHit{
		{ArticleId: 1, Score: 5},
		{ArticleId: 2, Score: 9},
		{ArticleId: 1, Score: 8},
		{ArticleId: 1, Score: 7},
		{ArticleId: 2, Score: 6},
		{ArticleId: 3, Score: 4},
	}

	got := dedupeChunkHits(hits, 2, 4)

	if len(got) != 4 {
		t.Fatalf("want 4 chunks, got %d", len(got))
	}
	wantOrder := []int64{2, 1, 1, 2}
	for i, want := range wantOrder {
		if got[i].ArticleId != want {
			t.Fatalf("chunk %d article = %d, want %d", i, got[i].ArticleId, want)
		}
	}
	if got[0].Score != 9 || got[3].Score != 6 {
		t.Fatalf("chunks should be sorted by score, got %v", scores(got))
	}
}

func TestDedupeChunkHitsRespectsLimit(t *testing.T) {
	hits := []*model.ArticleChunkHit{
		{ArticleId: 1, Score: 5},
		{ArticleId: 2, Score: 9},
		{ArticleId: 3, Score: 7},
	}

	got := dedupeChunkHits(hits, 2, 2)

	if len(got) != 2 || got[0].ArticleId != 2 || got[1].ArticleId != 3 {
		t.Fatalf("unexpected selection: %+v", got)
	}
}

func TestDedupeChunkHitsEmpty(t *testing.T) {
	if got := dedupeChunkHits(nil, 2, 5); len(got) != 0 {
		t.Fatalf("empty input should return nothing, got %d", len(got))
	}
	if got := dedupeChunkHits([]*model.ArticleChunkHit{{ArticleId: 1, Score: 1}}, 2, 0); len(got) != 0 {
		t.Fatalf("zero limit should return nothing, got %d", len(got))
	}
}

func TestTruncateRunes(t *testing.T) {
	tests := []struct {
		name string
		in   string
		max  int
		want string
	}{
		{name: "short", in: "博客", max: 5, want: "博客"},
		{name: "ascii", in: "abcdef", max: 3, want: "abc"},
		{name: "chinese", in: "博客助手测试", max: 2, want: "博客"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateRunes(tt.in, tt.max); got != tt.want {
				t.Fatalf("truncateRunes(%q, %d) = %q, want %q", tt.in, tt.max, got, tt.want)
			}
		})
	}
}

func scores(hits []*model.ArticleChunkHit) []float64 {
	out := make([]float64, 0, len(hits))
	for _, hit := range hits {
		out = append(out, hit.Score)
	}
	return out
}
