package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProxy struct {
	client openai.Client // v3 客户端实例
}

func NewOpenAIProxy(opts ...option.RequestOption) *OpenAIProxy {
	// 默认客户端配置
	defaultOpts := []option.RequestOption{
		option.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
	}
	// 合并自定义配置
	defaultOpts = append(defaultOpts, opts...)

	// 创建 OpenAI 客户端
	client := openai.NewClient(defaultOpts...)
	return &OpenAIProxy{
		client: client,
	}
}

// 通用错误响应方法（结构体内部方法，格式完全匹配 OpenAI）
func (p *OpenAIProxy) sendError(w http.ResponseWriter, statusCode int, message, errType string, param interface{}, code string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"type":    errType,
			"param":   param,
			"code":    code,
		},
	})
}

// ChatCompletionsHandler 处理 /v1/chat/completions 路由请求
func (p *OpenAIProxy) ChatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 读取并解析请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Failed to read request body: %v", err),
			"invalid_request_error", nil, "read_body_error")
		return
	}

	// 2. 先解析为 map 检查 stream 参数
	var rawReq map[string]interface{}
	if err = json.Unmarshal(body, &rawReq); err != nil {
		p.sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid request parameters: %v", err),
			"invalid_request_error", nil, "invalid_params")
		return
	}

	// 3. 判断是否为流式请求
	isStream := false
	if stream, ok := rawReq["stream"].(bool); ok && stream {
		isStream = true
	}

	// 4. 解析为 OpenAI 参数
	var reqParams openai.ChatCompletionNewParams
	if err = json.Unmarshal(body, &reqParams); err != nil {
		p.sendError(w, http.StatusBadRequest,
			fmt.Sprintf("Invalid request parameters: %v", err),
			"invalid_request_error", nil, "invalid_params")
		return
	}

	if isStream {
		p.handleStreamResponse(w, r, reqParams)
		return
	}

	// 5. 非流式响应
	resp, err := p.client.Chat.Completions.New(r.Context(), reqParams)
	if err != nil {
		var apiErr *openai.Error
		if errors.As(err, &apiErr) {
			p.sendError(w, apiErr.StatusCode, apiErr.Message, apiErr.Type, apiErr.Param, apiErr.Code)
			return
		}
		p.sendError(w, http.StatusInternalServerError,
			fmt.Sprintf("Proxy server error: %v", err),
			"server_error", nil, "proxy_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// handleStreamResponse 处理流式响应
func (p *OpenAIProxy) handleStreamResponse(w http.ResponseWriter, r *http.Request, reqParams openai.ChatCompletionNewParams) {
	// 设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	flusher, ok := w.(http.Flusher)
	if !ok {
		p.sendError(w, http.StatusInternalServerError, "Streaming not supported", "server_error", nil, "no_flusher")
		return
	}

	// 调用流式 API
	stream := p.client.Chat.Completions.NewStreaming(r.Context(), reqParams)

	// 逐块读取并发送
	for stream.Next() {
		chunk := stream.Current()
		data, _ := json.Marshal(chunk)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	if err := stream.Err(); err != nil {
		data, _ := json.Marshal(map[string]interface{}{
			"error": map[string]interface{}{
				"message": err.Error(),
				"type":    "stream_error",
			},
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	// 发送结束标记
	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}
