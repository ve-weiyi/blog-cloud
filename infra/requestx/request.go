package requestx

import (
	"net"
	"net/http"
	"strings"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
)

// ClientIP 取客户端真实 IP：优先反代覆写的 X-Real-IP，其次 X-Forwarded-For 的最后一跳，
// 最后回落到 RemoteAddr。
//
// 只取最后一跳：请求经反代时，只有最靠近网关的那一跳是反代写入的，其前缀由客户端
// 自行传入、不可信——按整串取值等于让一个请求头就能伪造 IP。
func ClientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get(bizheader.HeaderXRealIP)); ip != "" {
		return ip
	}

	if xff := r.Header.Get(bizheader.HeaderXForwardedFor); xff != "" {
		hops := strings.Split(xff, ",")
		if ip := strings.TrimSpace(hops[len(hops)-1]); ip != "" {
			return ip
		}
	}

	return extractIP(r.RemoteAddr)
}

func extractIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
