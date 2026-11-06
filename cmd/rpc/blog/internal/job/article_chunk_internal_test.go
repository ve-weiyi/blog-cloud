package job

import (
	"testing"
	"time"

	"github.com/ve-weiyi/blog-cloud/model"
)

func TestChunkCursorRoundTrip(t *testing.T) {
	cursor := chunkCursor{UpdatedAt: time.Unix(1700000000, 0), Id: 42}

	got := decodeChunkCursor(cursor.encode())

	if got.UpdatedAt.Unix() != cursor.UpdatedAt.Unix() || got.Id != cursor.Id {
		t.Fatalf("round trip mismatch: %+v != %+v", got, cursor)
	}
}

func TestDecodeChunkCursorInvalid(t *testing.T) {
	for _, raw := range []string{"", "abc", "1", "1|x", "x|1"} {
		got := decodeChunkCursor(raw)
		if !got.UpdatedAt.IsZero() || got.Id != 0 {
			t.Fatalf("%q should decode to zero cursor, got %+v", raw, got)
		}
	}
}

// 内容与标题都没变时跳过重写，任意一项变化都触发重建。
func TestSameChunks(t *testing.T) {
	current := []*model.TArticleChunk{
		{Ordinal: 0, ContentHash: "h1", ArticleTitle: "标题A"},
		{Ordinal: 1, ContentHash: "h2", ArticleTitle: "标题A"},
	}
	same := []*model.TArticleChunk{
		{Ordinal: 0, ContentHash: "h1", ArticleTitle: "标题A"},
		{Ordinal: 1, ContentHash: "h2", ArticleTitle: "标题A"},
	}
	if !sameChunks(current, same) {
		t.Fatal("identical chunks should skip rebuild")
	}

	changedHash := []*model.TArticleChunk{
		{Ordinal: 0, ContentHash: "h1", ArticleTitle: "标题A"},
		{Ordinal: 1, ContentHash: "h3", ArticleTitle: "标题A"},
	}
	if sameChunks(current, changedHash) {
		t.Fatal("hash change should trigger rebuild")
	}

	changedTitle := []*model.TArticleChunk{
		{Ordinal: 0, ContentHash: "h1", ArticleTitle: "标题B"},
		{Ordinal: 1, ContentHash: "h2", ArticleTitle: "标题B"},
	}
	if sameChunks(current, changedTitle) {
		t.Fatal("title change should trigger rebuild")
	}

	if sameChunks(current, current[:1]) {
		t.Fatal("chunk count change should trigger rebuild")
	}
}
