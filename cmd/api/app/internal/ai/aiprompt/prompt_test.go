package aiprompt

import (
	"strings"
	"testing"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/contentservice"
	"github.com/ve-weiyi/blog-cloud/infra/ai"
)

func TestBuildHistoryKeepsRecentTurns(t *testing.T) {
	history := []ai.Message{
		{Role: "user", Content: "q1"}, {Role: "assistant", Content: "a1"},
		{Role: "user", Content: "q2"}, {Role: "assistant", Content: "a2"},
		{Role: "user", Content: "q3"}, {Role: "assistant", Content: "a3"},
	}

	got := BuildHistory(history, 2, 0)

	if len(got) != 4 {
		t.Fatalf("want last 2 turns (4 messages), got %d", len(got))
	}
	if got[0].Content != "q2" || got[3].Content != "a3" {
		t.Fatalf("unexpected window: %+v", got)
	}
}

func TestBuildHistoryDropsOldestByChars(t *testing.T) {
	history := []ai.Message{
		{Role: "user", Content: strings.Repeat("a", 10)},
		{Role: "assistant", Content: strings.Repeat("b", 10)},
		{Role: "user", Content: strings.Repeat("c", 10)},
	}

	got := BuildHistory(history, 0, 15)

	if len(got) != 1 || got[0].Content != strings.Repeat("c", 10) {
		t.Fatalf("unexpected history: %+v", got)
	}
}

func TestReferenceBlockDedupesCitations(t *testing.T) {
	chunks := []*contentservice.ArticleChunk{
		{ArticleId: 1, ArticleTitle: "文章A", Heading: "第一节", Content: "内容一"},
		{ArticleId: 1, ArticleTitle: "文章A", Heading: "第二节", Content: "内容二"},
		{ArticleId: 2, ArticleTitle: "文章B", Heading: "", Content: "内容三"},
	}

	block, citations := ReferenceBlock(chunks, 10000)

	if len(citations) != 2 || citations[0].ArticleId != 1 || citations[1].ArticleId != 2 {
		t.Fatalf("citations = %+v", citations)
	}
	for _, want := range []string{"[1]", "[2]", "[3]", "文章A", "正文", "内容三"} {
		if !strings.Contains(block, want) {
			t.Fatalf("block missing %q: %s", want, block)
		}
	}
}

func TestReferenceBlockRespectsCharLimit(t *testing.T) {
	chunks := []*contentservice.ArticleChunk{
		{ArticleId: 1, ArticleTitle: "A", Heading: "H", Content: strings.Repeat("字", 100)},
	}

	block, citations := ReferenceBlock(chunks, 10)

	if block != "" || len(citations) != 0 {
		t.Fatalf("oversized chunk should be dropped, got %q / %+v", block, citations)
	}
}

func TestBuildMessagesTrimsHistoryFirst(t *testing.T) {
	history := []ai.Message{
		{Role: "user", Content: strings.Repeat("旧", 200)},
		{Role: "assistant", Content: strings.Repeat("旧", 200)},
	}

	messages := BuildMessages("", history, "当前问题", 300)

	if len(messages) != 2 {
		t.Fatalf("want system + question only, got %d", len(messages))
	}
	if messages[0].Role != "system" || messages[1].Content != "当前问题" {
		t.Fatalf("unexpected messages: %+v", messages)
	}
}

func TestBuildMessagesKeepsHistoryWhenFits(t *testing.T) {
	history := []ai.Message{{Role: "user", Content: "上次的问题"}}

	messages := BuildMessages("<references>独特资料标记</references>", history, "当前问题", 10000)

	if len(messages) != 4 {
		t.Fatalf("want system + references + history + question, got %d", len(messages))
	}
	if messages[0].Role != "system" {
		t.Fatalf("first message should be system, got %q", messages[0].Role)
	}
	if strings.Contains(messages[0].Content, "独特资料标记") {
		t.Fatalf("references should not live in system: %q", messages[0].Content)
	}
	if messages[1].Role != "user" || !strings.Contains(messages[1].Content, "独特资料标记") {
		t.Fatalf("references should be a separate user message: %+v", messages[1])
	}
	if messages[3].Content != "当前问题" {
		t.Fatalf("last message should be the question, got %+v", messages[3])
	}
}

func TestErrorEvent(t *testing.T) {
	tests := []struct {
		err  error
		code string
	}{
		{err: &ai.Error{Kind: ai.KindTimeout}, code: "timeout"},
		{err: &ai.Error{Kind: ai.KindRateLimited}, code: "rate_limited"},
		{err: &ai.Error{Kind: ai.KindAuth}, code: "unavailable"},
		{err: &ai.Error{Kind: ai.KindUpstream}, code: "upstream_error"},
		{err: &ai.Error{Kind: ai.KindCanceled}, code: ""},
	}

	for _, tt := range tests {
		event := ErrorEvent(tt.err)
		if tt.code == "" {
			if event != nil {
				t.Fatalf("canceled should not emit an event, got %+v", event)
			}
			continue
		}
		if event == nil || event.Code != tt.code {
			t.Fatalf("err %v => %+v, want code %s", tt.err, event, tt.code)
		}
	}
}
