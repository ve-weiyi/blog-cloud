package requestx

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestClientIP 守住客户端真实 IP 的口径：只信反代写入的转发头，且只取最后一跳。
// 客户端自行传入的转发头前缀不可信——否则一个请求头即可伪造 IP，绕过 IP 维度的限流。
func TestClientIP(t *testing.T) {
	cases := []struct {
		name       string
		remoteAddr string
		realIP     string
		xff        string
		want       string
	}{
		{"优先取反代覆写的 X-Real-IP", "10.0.0.1:5555", "1.2.3.4", "", "1.2.3.4"},
		{"X-Real-IP 空白时回落到 XFF", "10.0.0.1:5555", "   ", "1.1.1.1, 2.2.2.2", "2.2.2.2"},
		{"只有 XFF 时取最后一跳", "10.0.0.1:5555", "", "1.1.1.1, 2.2.2.2, 3.3.3.3", "3.3.3.3"},
		{"XFF 单值时取该值", "10.0.0.1:5555", "", "1.1.1.1", "1.1.1.1"},
		{"XFF 带空白时去除空白", "10.0.0.1:5555", "", "1.1.1.1 ,  2.2.2.2 ", "2.2.2.2"},
		{"无转发头时回落到 RemoteAddr", "10.0.0.1:5555", "", "", "10.0.0.1"},
		{"转发头为空白时回落到 RemoteAddr", "10.0.0.1:5555", "", "   ", "10.0.0.1"},
		{"RemoteAddr 不带端口时原样返回", "10.0.0.1", "", "", "10.0.0.1"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/v1/article", nil)
			r.RemoteAddr = c.remoteAddr
			if c.realIP != "" {
				r.Header.Set("X-Real-IP", c.realIP)
			}
			if c.xff != "" {
				r.Header.Set("X-Forwarded-For", c.xff)
			}

			if got := ClientIP(r); got != c.want {
				t.Errorf("ClientIP = %q, want %q", got, c.want)
			}
		})
	}
}
