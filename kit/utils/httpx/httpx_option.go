package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client 表示HTTP客户端。
type Client struct {
	timeout time.Duration // 请求超时时间

	method string // 请求方法
	url    string // 请求URL

	headers map[string]string // 请求头部
	params  map[string]string // 请求路径上的查询参数
	body    []byte            // 请求体
}

// Option 表示用于配置HTTP客户端的函数选项。
type Option func(*Client)

// NewClient 使用默认设置创建一个新的HTTP客户端。
func NewClient(method, url string, options ...Option) *Client {
	client := &Client{
		method: method,
		url:    url,

		timeout: 30 * time.Second, // 默认超时时间
		headers: make(map[string]string),
		params:  make(map[string]string),
		body:    make([]byte, 0),
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
func (c *Client) DoRequest() (respBody []byte, err error) {
	// 解析URL
	uv, err := url.ParseRequestURI(c.url)
	if err != nil {
		return nil, err
	}

	// 使用请求体创建请求
	req, err := http.NewRequest(c.method, uv.String(), bytes.NewBuffer(c.body))
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

	// 底层HTTP客户端
	httpClient := &http.Client{}

	// 设置超时时间
	httpClient.Timeout = c.timeout

	// 执行请求
	resp, err := httpClient.Do(req)
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
func (c *Client) EncodeURL() string {
	req, err := http.NewRequest(c.method, c.url, nil)
	if err != nil {
		return c.url
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
func (c *Client) CURL() string {
	var curl string
	curl = fmt.Sprintf("curl --location --request %s '%s'", c.method, c.url)

	for key, value := range c.headers {
		curl += " \\\n"
		curl += fmt.Sprintf("--header '%s: %s'", key, value)
	}

	if len(c.body) > 0 {
		curl += " \\\n"
		curl += fmt.Sprintf("--data-raw '%s'", string(c.body))
	}

	return curl
}
