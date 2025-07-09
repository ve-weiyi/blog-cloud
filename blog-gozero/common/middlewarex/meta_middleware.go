package middlewarex

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

type CtxMetaMiddleware struct {
}

func NewCtxMetaMiddleware() *CtxMetaMiddleware {
	return &CtxMetaMiddleware{}
}

// 将http header 放入 ctx 里面使用 metadata 保存.
func (m *CtxMetaMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("CtxMetaMiddleware Handle")
		ctx := r.Context()
		md := metadata.MD{}
		// 从真实http 头获取header
		if len(r.Header) > 0 {
			for k, v := range r.Header {
				if len(k) == 0 {
					continue
				}
				var value string
				if len(v) > 0 {
					value = v[0]
				}

				keyLowercase := strings.ToLower(k)
				//logx.Infof("add k=%s, v=%+v", keyLowercase, value)
				for _, key := range restx.HeaderFields {
					if key == keyLowercase {
						md.Set(key, value)
						ctx = context.WithValue(ctx, key, value)
					}
				}
			}
		}

		md.Set(restx.HeaderRPCRemoteAgent, r.UserAgent())
		md.Set(restx.HeaderRPCRemoteIP, httpx.GetRemoteAddr(r))

		ctx = metadata.NewOutgoingContext(ctx, md)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
