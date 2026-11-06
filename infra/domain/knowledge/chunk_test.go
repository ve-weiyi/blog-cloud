package knowledge

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSplitHeadingPath(t *testing.T) {
	markdown := "# 一级标题\n\n正文一\n\n## 二级标题\n\n正文二\n\n# 另一个一级\n\n正文三"
	chunks := Split(markdown)

	if len(chunks) != 3 {
		t.Fatalf("want 3 chunks, got %d", len(chunks))
	}
	if chunks[0].Heading != "一级标题" || chunks[1].Heading != "一级标题 / 二级标题" {
		t.Fatalf("unexpected heading path: %q, %q", chunks[0].Heading, chunks[1].Heading)
	}
	if chunks[2].Heading != "另一个一级" {
		t.Fatalf("heading stack should pop closed levels, got %q", chunks[2].Heading)
	}
	for i, chunk := range chunks {
		if chunk.Ordinal != i {
			t.Fatalf("ordinal = %d, want %d", chunk.Ordinal, i)
		}
	}
}

// 代码块里的 # 不是标题。
func TestSplitIgnoresFencedHeadings(t *testing.T) {
	markdown := "# 真标题\n\n```bash\n# 这是注释\n\necho hi\n```\n\n正文"
	chunks := Split(markdown)

	if len(chunks) != 1 {
		t.Fatalf("want 1 chunk, got %d", len(chunks))
	}
	if chunks[0].Heading != "真标题" {
		t.Fatalf("unexpected heading %q", chunks[0].Heading)
	}
	if !strings.Contains(chunks[0].Content, "# 这是注释") {
		t.Fatalf("fenced content should stay in chunk: %q", chunks[0].Content)
	}
}

// 超长段落必须切分，且不丢内容。
func TestSplitLongParagraph(t *testing.T) {
	body := strings.Repeat("好", 2000)
	chunks := Split("# 标题\n\n" + body)

	if len(chunks) < 3 {
		t.Fatalf("want at least 3 chunks, got %d", len(chunks))
	}
	total := 0
	for _, chunk := range chunks {
		n := utf8.RuneCountInString(chunk.Content)
		if n > MaxContentRunes {
			t.Fatalf("chunk too long: %d", n)
		}
		total += n
	}
	if total != 2000 {
		t.Fatalf("content lost while splitting: %d != 2000", total)
	}
}

func TestSplitKeepsSentenceBoundary(t *testing.T) {
	first := strings.Repeat("甲", MaxContentRunes-10) + "。"
	body := first + strings.Repeat("乙", 100)
	chunks := Split(body)

	if len(chunks) < 2 {
		t.Fatalf("want at least 2 chunks, got %d", len(chunks))
	}
	if !strings.HasSuffix(chunks[0].Content, "。") {
		t.Fatalf("should cut after sentence mark, got tail %q", lastRunes(chunks[0].Content, 5))
	}
}

func TestSplitEmpty(t *testing.T) {
	if chunks := Split("   \n\n"); len(chunks) != 0 {
		t.Fatalf("empty markdown should produce no chunks, got %d", len(chunks))
	}
}

func TestSplitSummaryAndHash(t *testing.T) {
	content := strings.Repeat("长", SummaryRunes+50)
	first := Split(content)
	second := Split(content)

	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("want single chunk, got %d / %d", len(first), len(second))
	}
	if utf8.RuneCountInString(first[0].Summary) > SummaryRunes+1 {
		t.Fatalf("summary too long: %d", utf8.RuneCountInString(first[0].Summary))
	}
	if first[0].Hash != second[0].Hash {
		t.Fatal("hash should be deterministic")
	}
	if first[0].Summary == first[0].Content {
		t.Fatal("summary should be truncated")
	}
}

func lastRunes(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[len(runes)-n:])
}
