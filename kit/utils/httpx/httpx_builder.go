package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ClientBuilder 用于构建 HTTP 客户端。
type ClientBuilder struct {
	method string // 请求方法
	url    string // 请求URL

	timeout time.Duration     // 请求超时时间
	headers map[string]string // 请求头部
	params  map[string]string // 请求参数
	body    []byte            // 请求体
}

// NewClientBuilder 创建一个具有默认设置的新 ClientBuilder。
func NewClientBuilder(method, url string) *ClientBuilder {
	client := &ClientBuilder{
		method: method,
		url:    url,

		timeout: 30 * time.Second, // 默认超时时间
		headers: make(map[string]string),
		params:  make(map[string]string),
		body:    make([]byte, 0),
	}

	return client
}

// WithTimeout 设置 HTTP 请求的超时时间。
func (c *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
	c.timeout = timeout
	return c
}

// WithHeaders 设置 HTTP 请求的头部。
func (c *ClientBuilder) WithHeaders(headers map[string]string) *ClientBuilder {
	c.headers = headers
	return c
}

// WithParams 设置 HTTP 请求的参数。
func (c *ClientBuilder) WithParams(params map[string]string) *ClientBuilder {
	c.params = params
	return c
}

// WithBodyJson 设置 HTTP 请求的正文。
func (c *ClientBuilder) WithBody(body []byte) *ClientBuilder {
	c.body = body
	return c
}

// Build 根据构建器的设置创建一个新的 HTTP 客户端。
func (c *ClientBuilder) Build() *http.Client {
	client := &http.Client{
		Timeout: c.timeout,
	}

	return client
}

// DoRequest 执行一个 HTTP 请求。
func (c *ClientBuilder) DoRequest() (respBody []byte, err error) {
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
