package middlewarex

import (
	"net/http"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin") //请求头部
		if origin != "" {
			// 允许当前源
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE, UPDATE")
			// 允许跨域设置可以返回其他子段
			w.Header().Set("Access-Control-Allow-Headers", "*")
			// 允许浏览器（客户端）可以解析的头部
			w.Header().Set("Access-Control-Expose-Headers", "*")
			// 允许携带cookie
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			// 预检请求缓存 24 小时
			w.Header().Set("Access-Control-Max-Age", "7200")
		} else {
			// 允许所有源
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// 处理 OPTIONS 预检请求：直接返回 200，不执行后续逻辑
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 执行下一个处理器（核心业务逻辑）
		next(w, r)
	}
}
