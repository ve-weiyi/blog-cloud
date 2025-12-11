package httpx

import (
	"net/http"
	"time"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   30 * time.Second, // 整个请求的超时
	}
}
