package middlewarex

import (
	"context"
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"

	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/http/requestx"
)

// clientIdentityHeaders 是客户端声明身份的头，不能直接进入 RPC metadata：
// 身份类由认证中间件在鉴权后写入可信值，平台标识由网关按自身配置写入。
var clientIdentityHeaders = map[string]struct{}{
	headerx.HeaderUid:           {},
	headerx.HeaderAuthorization: {},
	headerx.HeaderToken:         {},
	headerx.HeaderXAdminId:      {},
	headerx.HeaderXAdminAuth:    {},
	headerx.HeaderAppName:       {},
}

type CtxMetaMiddleware struct {
	// appName 是本网关的平台标识；为空时不写入该头
	appName string
}

func NewCtxMetaMiddleware(appName string) *CtxMetaMiddleware {
	return &CtxMetaMiddleware{appName: appName}
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
				for _, key := range headerx.HeaderFields {
					if key == keyLowercase {
						if _, ok := clientIdentityHeaders[key]; ok {
							continue
						}
						md.Set(key, value)
						ctx = context.WithValue(ctx, key, value)
					}
				}
			}
		}

		// 平台标识取自网关配置，不用客户端声明的值：同一个人从两个端进来时，
		// 「这是哪个端发起的请求」只能由网关盖章。
		if m.appName != "" {
			md.Set(headerx.HeaderAppName, m.appName)
			ctx = context.WithValue(ctx, headerx.HeaderAppName, m.appName)
		}

		md.Set(metax.HeaderRPCRemoteAgent, r.UserAgent())
		md.Set(metax.HeaderRPCRemoteIP, requestx.ClientIP(r))

		ctx = metadata.NewOutgoingContext(ctx, md)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
