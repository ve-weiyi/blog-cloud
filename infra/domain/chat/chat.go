// Package chat 定义聊天内容的公共约束，供 RPC 与网关钩子共用。
package chat

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	// MaxContentRunes 单条消息的最大字符数。客户端也会限制，这里是权威边界。
	MaxContentRunes = 1000
	TypeText        = "text"
)

// NormalizeContent 去空白并校验长度；返回写入数据库的文本。
func NormalizeContent(raw string) (string, error) {
	content := strings.TrimSpace(raw)
	n := utf8.RuneCountInString(content)
	if n == 0 || n > MaxContentRunes {
		return "", fmt.Errorf("chat: content must be 1-%d characters", MaxContentRunes)
	}
	return content, nil
}

// ValidType 判断消息类型是否受支持；当前只开放文本。
func ValidType(messageType string) bool {
	return messageType == TypeText
}
