package httpx

import (
	"encoding/json"
)

// WithHeader 设置请求头
func WithHeader(key string, value string) Option {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithParam 设置查询参数
func WithParam(key string, value string) Option {
	return func(c *Client) {
		c.params[key] = value
	}
}

// WithBodyObject 设置请求体
func WithBodyObject(obj interface{}) Option {
	return func(c *Client) {
		data, err := json.Marshal(obj)
		if err != nil {
			return
		}

		c.body = data
	}
}

func structToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func jsonToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
