package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client 表示HTTP客户端。
type Client struct {
	httpClient *http.Client  // 底层HTTP客户端
	timeout    time.Duration // 请求超时时间

	headers map[string]string // 请求头部
	params  map[string]string // 请求参数
	body    []byte            // 请求体
}

// Option 表示用于配置HTTP客户端的函数选项。
type Option func(*Client)

// NewClient 使用默认设置创建一个新的HTTP客户端。
func NewClient(options ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{},
		timeout:    30 * time.Second, // 默认超时时间
		headers:    make(map[string]string),
		params:     make(map[string]string),
	}

	// 应用选项
	for _, option := range options {
		option(client)
	}

	return client
}

// WithTimeout 设置HTTP请求的超时时间。
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithHeaders 设置HTTP请求的头部。
func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

// WithParams 设置HTTP请求的参数。
func WithParams(params map[string]string) Option {
	return func(c *Client) {
		for key, value := range params {
			c.params[key] = value
		}
	}
}

// WithBody 设置HTTP请求的请求体。
func WithBody(body []byte) Option {
	return func(c *Client) {
		c.body = body
	}
}

// DoRequest 执行一个HTTP请求。
func (c *Client) DoRequest(method, url string) (respBody []byte, err error) {
	// 使用请求体创建请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(c.body))
	if err != nil {
		return nil, err
	}

	// 设置头部
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// 设置查询参数
	query := req.URL.Query()
	for key, value := range c.params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	// 设置超时时间
	c.httpClient.Timeout = c.timeout

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request fail. url:%v, code:%d,err:%s", req.URL.String(), resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// EncodeURL 编码URL。
func (c *Client) EncodeURL(rawURL string) string {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return rawURL
	}

	// 设置查询参数
	query := req.URL.Query()
	for key, value := range c.params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()

	return req.URL.String()
}

// 输出 curl
func (c *Client) CURL(method, url string) (string, error) {
	var curl string
	curl = fmt.Sprintf("curl -X %s %s", method, url)

	for key, value := range c.headers {
		curl += fmt.Sprintf(" -H '%s: %s'", key, value)
	}

	if len(c.body) > 0 {
		curl += fmt.Sprintf(" -d '%s'", string(c.body))
	}

	return curl, nil
}
