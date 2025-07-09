package requestx

import (
	"context"
	"net/http"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
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
		Token:     r.Header.Get(restx.HeaderToken),
		Uid:       r.Header.Get(restx.HeaderUid),
		IpAddress: r.RemoteAddr,
		UserAgent: r.UserAgent(),
	}

	return reqCtx
}
