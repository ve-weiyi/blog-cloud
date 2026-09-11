package requestx

import (
	"context"
	"net"
	"net/http"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
)

// 请求上下文,一般存放请求头参数
type Context struct {
	context.Context `json:"-" header:"-"`
	Token           string `json:"token" header:"token" example:""`
	Uid             string `json:"uid" header:"-" example:""`
	IpAddress       string `json:"ip_address" header:"-" example:""`
	UserAgent       string `json:"user_agent" header:"-" example:""`
}

func (s *Context) GetContext() context.Context {
	return s.Context
}

// 获取请求上下文
func ParseRequestContext(r *http.Request) *Context {
	reqCtx := &Context{
		Context:   r.Context(),
		Token:     r.Header.Get(bizheader.HeaderToken),
		Uid:       r.Header.Get(bizheader.HeaderUid),
		IpAddress: extractIP(r.RemoteAddr),
		UserAgent: r.UserAgent(),
	}

	return reqCtx
}

func extractIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
