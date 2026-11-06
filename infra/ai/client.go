// Package ai 是 OpenAI 兼容流式接口的最小客户端，只覆盖对话生成所需的能力。
package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// KindAuth 上游拒绝凭据。
	KindAuth = "upstream_auth"
	// KindRateLimited 上游限流。
	KindRateLimited = "upstream_rate_limited"
	// KindTimeout 上游超时。
	KindTimeout = "upstream_timeout"
	// KindUpstream 上游返回其他错误。
	KindUpstream = "upstream_error"
	// KindCanceled 调用方取消。
	KindCanceled = "canceled"
)

// Message 是一条对话消息。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Usage 是上游报告的 token 用量。
type Usage struct {
	PromptTokens     int64
	CompletionTokens int64
}

// Config 是客户端配置。
type Config struct {
	BaseURL     string
	APIKey      string
	Model       string
	Temperature float64
	MaxTokens   int
	Timeout     time.Duration
}

// Error 是上游错误，Kind 用于上层选择展示文案。
type Error struct {
	Kind    string
	Status  int
	Message string
}

func (e *Error) Error() string {
	if e.Status != 0 {
		return fmt.Sprintf("ai: %s (status %d): %s", e.Kind, e.Status, e.Message)
	}
	return fmt.Sprintf("ai: %s: %s", e.Kind, e.Message)
}

// Client 走 HTTP 的 chat/completions 流式接口。
type Client struct {
	cfg        Config
	httpClient *http.Client
}

// NewClient 创建客户端。
//
// Timeout 是整条流的总时限（含响应体读取），不是空闲超时：超过它会被判为
// 超时并中断生成，因此配置值必须覆盖一次完整回答的耗时。
func NewClient(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: timeout},
	}
}

// Model 返回当前配置的模型名。
func (c *Client) Model() string {
	return c.cfg.Model
}

// Stream 发起流式请求，每收到一段正文调用一次 onDelta。
//
// 返回值是上游报告的用量；onDelta 返回错误或 ctx 取消时立即停止读取。
func (c *Client) Stream(ctx context.Context, messages []Message, onDelta func(string) error) (Usage, error) {
	var usage Usage

	body, err := c.encode(messages)
	if err != nil {
		return usage, err
	}

	url := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return usage, fmt.Errorf("ai: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return usage, wrapTransportError(ctx, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return usage, decodeError(resp)
	}

	scanner := bufio.NewScanner(resp.Body)
	// 单条事件可能包含较长的正文片段，放宽单行上限
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}

		var chunk chatChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			continue
		}
		if chunk.Usage != nil {
			usage.PromptTokens = chunk.Usage.PromptTokens
			usage.CompletionTokens = chunk.Usage.CompletionTokens
		}
		for _, choice := range chunk.Choices {
			if choice.Delta.Content == "" {
				continue
			}
			if err := onDelta(choice.Delta.Content); err != nil {
				// 客户端断开时 onDelta 通常返回 ctx 的错误，统一归到取消而不是上游失败
				if ctx.Err() != nil {
					return usage, wrapTransportError(ctx, err)
				}
				return usage, err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return usage, wrapTransportError(ctx, err)
	}
	return usage, nil
}

func (c *Client) encode(messages []Message) ([]byte, error) {
	payload := chatRequest{
		Model:       c.cfg.Model,
		Messages:    messages,
		Stream:      true,
		Temperature: c.cfg.Temperature,
		MaxTokens:   c.cfg.MaxTokens,
		StreamOptions: &streamOptions{
			IncludeUsage: true,
		},
	}
	return json.Marshal(payload)
}

type chatRequest struct {
	Model         string         `json:"model"`
	Messages      []Message      `json:"messages"`
	Stream        bool           `json:"stream"`
	Temperature   float64        `json:"temperature,omitempty"`
	MaxTokens     int            `json:"max_tokens,omitempty"`
	StreamOptions *streamOptions `json:"stream_options,omitempty"`
}

type streamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type chatChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	Usage *struct {
		PromptTokens     int64 `json:"prompt_tokens"`
		CompletionTokens int64 `json:"completion_tokens"`
	} `json:"usage"`
}

func decodeError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	kind := KindUpstream
	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		kind = KindAuth
	case resp.StatusCode == http.StatusTooManyRequests:
		kind = KindRateLimited
	}
	return &Error{Kind: kind, Status: resp.StatusCode, Message: message}
}

func wrapTransportError(ctx context.Context, err error) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return &Error{Kind: KindCanceled, Message: "request canceled"}
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(err, context.DeadlineExceeded) {
		return &Error{Kind: KindTimeout, Message: "upstream timeout"}
	}

	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) && netErr.Timeout() {
		return &Error{Kind: KindTimeout, Message: err.Error()}
	}
	return &Error{Kind: KindUpstream, Message: err.Error()}
}
