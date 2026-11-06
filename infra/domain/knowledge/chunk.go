// Package knowledge 把文章 Markdown 切成分块，供全文检索使用。
package knowledge

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

const (
	// MaxContentRunes 单个分块的最大字符数。
	MaxContentRunes = 800
	// SummaryRunes 分块摘要的最大字符数，超出的直接截断。
	SummaryRunes = 160
	// MaxHeadingRunes 标题路径的最大字符数，与 heading 列宽一致。
	MaxHeadingRunes = 255
)

var (
	headingPattern = regexp.MustCompile(`^(#{1,6})[ \t]+(.+?)[ \t#]*$`)
	// sentenceBreak 是硬切长段落时优先选择的中文/英文句末符号。
	sentenceBreak = []rune("。！？!?；;\n")
)

// Chunk 是一个检索单元。
type Chunk struct {
	Ordinal int    // 在文章内的序号
	Heading string // 标题路径，如 "一级 / 二级"
	Content string // 正文
	Summary string // 截断摘要
	Hash    string // 内容摘要，用于判断是否需要重建
}

type heading struct {
	level int
	text  string
}

// Split 按标题切段、按段落聚合，把 Markdown 切成不超过 MaxContentRunes 的分块。
//
// 代码块里的 # 不算标题；标题路径只保留当前层级的祖先链。
func Split(markdown string) []Chunk {
	markdown = strings.ReplaceAll(markdown, "\r\n", "\n")
	lines := strings.Split(markdown, "\n")

	var (
		chunks  []Chunk
		stack   []heading
		section []string
		inFence bool
		fence   string
	)

	flush := func() {
		text := strings.TrimSpace(strings.Join(section, "\n"))
		section = nil
		if text == "" {
			return
		}
		path := headingPath(stack)
		for _, part := range splitSection(text) {
			chunks = append(chunks, newChunk(len(chunks), path, part))
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if !inFence {
			if marker, ok := fenceMarker(trimmed); ok {
				inFence = true
				fence = marker
				section = append(section, line)
				continue
			}
		} else {
			section = append(section, line)
			if strings.HasPrefix(trimmed, fence) {
				inFence = false
			}
			continue
		}

		if match := headingPattern.FindStringSubmatch(trimmed); match != nil {
			flush()
			level := len(match[1])
			for len(stack) > 0 && stack[len(stack)-1].level >= level {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, heading{level: level, text: strings.TrimSpace(match[2])})
			continue
		}

		section = append(section, line)
	}

	flush()
	return chunks
}

func newChunk(ordinal int, path, content string) Chunk {
	path = truncateRunes(path, MaxHeadingRunes)
	content = strings.TrimSpace(content)
	sum := sha256.Sum256([]byte(path + "\n" + content))
	return Chunk{
		Ordinal: ordinal,
		Heading: path,
		Content: content,
		Summary: summarize(content),
		Hash:    hex.EncodeToString(sum[:]),
	}
}

func headingPath(stack []heading) string {
	parts := make([]string, 0, len(stack))
	for _, h := range stack {
		parts = append(parts, h.text)
	}
	return strings.Join(parts, " / ")
}

func fenceMarker(line string) (string, bool) {
	if strings.HasPrefix(line, "```") {
		return "```", true
	}
	if strings.HasPrefix(line, "~~~") {
		return "~~~", true
	}
	return "", false
}

// splitSection 先把整节按空行切成段落，再把段落聚合到上限内。
func splitSection(text string) []string {
	var (
		parts     []string
		paragraph []string
		flushPara = func() {
			if len(paragraph) == 0 {
				return
			}
			parts = append(parts, strings.TrimSpace(strings.Join(paragraph, "\n")))
			paragraph = nil
		}
		result      []string
		buffer      []rune
		flushBuffer = func() {
			if len(buffer) == 0 {
				return
			}
			result = append(result, strings.TrimSpace(string(buffer)))
			buffer = nil
		}
	)

	for _, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			flushPara()
			continue
		}
		paragraph = append(paragraph, line)
	}
	flushPara()

	for _, part := range parts {
		runes := []rune(part)
		for len(runes) > 0 {
			// 剩余段落还能整段放下
			if len(buffer) == 0 && len(runes) <= MaxContentRunes {
				buffer = append(buffer, runes...)
				break
			}

			space := MaxContentRunes - len(buffer) - 2 // 预留段落之间的空行
			if space <= 0 {
				flushBuffer()
				continue
			}
			if len(runes) <= space {
				if len(buffer) > 0 {
					buffer = append(buffer, '\n', '\n')
				}
				buffer = append(buffer, runes...)
				break
			}

			cut := sentenceCut(runes, space)
			if len(buffer) > 0 {
				buffer = append(buffer, '\n', '\n')
			}
			buffer = append(buffer, runes[:cut]...)
			flushBuffer()
			runes = runes[cut:]
			runes = []rune(strings.TrimLeft(string(runes), "\n "))
		}
	}
	flushBuffer()

	return result
}

// sentenceCut 在 [space*3/5, space] 内就近找句末符号，找不到就直接在 space 处切。
func sentenceCut(runes []rune, space int) int {
	if space >= len(runes) {
		return len(runes)
	}
	lower := space * 3 / 5
	for i := space - 1; i >= lower; i-- {
		for _, mark := range sentenceBreak {
			if runes[i] == mark {
				return i + 1
			}
		}
	}
	return space
}

func truncateRunes(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max])
}

func summarize(content string) string {
	runes := []rune(strings.Join(strings.Fields(content), " "))
	if len(runes) <= SummaryRunes {
		return string(runes)
	}
	return string(runes[:SummaryRunes]) + "…"
}
