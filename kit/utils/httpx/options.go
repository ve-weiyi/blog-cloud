package httpx

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
)

// WithHeader 设置请求头
func WithHeader(key string, value string) Option {
	return func(c *Request) {
		c.headers[key] = value
	}
}

// WithParam 设置查询参数
func WithParam(key string, value string) Option {
	return func(c *Request) {
		c.params[key] = value
	}
}

// WithBodyJson 设置请求体
func WithBodyJson(obj interface{}) Option {
	return func(c *Request) {
		data, err := json.Marshal(obj)
		if err != nil {
			return
		}

		c.headers["Content-Type"] = "application/json"
		c.body = data
	}
}

// WithBodyForm 设置请求体
func WithBodyForm(obj map[string]string) Option {
	return func(c *Request) {
		data := url.Values{}
		for key, value := range obj {
			data.Set(key, value)
		}

		c.headers["Content-Type"] = "application/x-www-form-urlencoded"
		c.body = []byte(data.Encode())
	}
}

// WithBasicAuth 设置请求体
func WithBasicAuth(username string, password string) Option {
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
	return func(c *Request) {
		c.headers["Authorization"] = fmt.Sprintf("Basic %s", auth)
	}
}
