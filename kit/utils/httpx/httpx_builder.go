package httpx

// ClientBuilder 用于构建 HTTP 客户端。
type ClientBuilder struct {
	method string // 请求方法
	url    string // 请求URL

	headers map[string]string // 请求头部
	params  map[string]string // 请求参数
	body    []byte            // 请求体
}

// NewClientBuilder 创建一个具有默认设置的新 ClientBuilder。
func NewClientBuilder(method, url string) *ClientBuilder {
	client := &ClientBuilder{
		method: method,
		url:    url,

		headers: make(map[string]string),
		params:  make(map[string]string),
		body:    make([]byte, 0),
	}

	return client
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
func (c *ClientBuilder) Build() *Request {
	return &Request{
		method:  c.method,
		url:     c.url,
		headers: c.headers,
		params:  c.params,
		body:    c.body,
	}
}
