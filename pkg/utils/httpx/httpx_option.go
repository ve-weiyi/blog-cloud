package httpx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request 表示HTTP客户端。
type Request struct {
	method string // 请求方法
	url    string // 请求URL

	headers map[string]string // 请求头部
	params  map[string]string // 请求路径上的查询参数
	body    []byte            // 请求体
}

// Option 表示用于配置HTTP客户端的函数选项。
type Option func(*Request)

// NewRequest 使用默认设置创建一个新的HTTP客户端。
func NewRequest(method, url string, options ...Option) *Request {
	client := &Request{
		method: method,
		url:    url,

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

// WithHeaders 设置HTTP请求的头部。
func WithHeaders(headers map[string]string) Option {
	return func(c *Request) {
		for key, value := range headers {
			c.headers[key] = value
		}
	}
}

// WithParams 设置HTTP请求的参数。
func WithParams(params map[string]string) Option {
	return func(c *Request) {
		for key, value := range params {
			c.params[key] = value
		}
	}
}

// WithBody 设置HTTP请求的请求体。
func WithBody(body []byte) Option {
	return func(c *Request) {
		c.body = body
	}
}

// Do 执行一个HTTP请求。
func (c *Request) Do() (respBody []byte, err error) {
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
	//httpClient := &http.Client{}

	// 执行请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(`http response abnormal status: %s`, resp.Status)
	}

	defer resp.Body.Close()
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

// EncodeURL 编码URL。
func (c *Request) EncodeURL() string {
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
func (c *Request) CURL() string {
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
