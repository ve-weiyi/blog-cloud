package middleware

import (
	"net/http"

	"github.com/go-openapi/spec"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
)

type VisitLogMiddleware struct {
	sp *spec.Swagger

	sr syslogrpc.SyslogRpc
}

func NewVisitLogMiddleware(sp *spec.Swagger, sr syslogrpc.SyslogRpc) *VisitLogMiddleware {
	return &VisitLogMiddleware{
		sp: sp,
		sr: sr,
	}
}

// 记录操作记录
func (m *VisitLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("VisitLogMiddleware Handle path: %v", r.URL.Path)

		// 调用下一层的处理
		next.ServeHTTP(w, r)

		var module string
		api := m.getApi(r.URL.Path, r.Method)
		if api != nil {
			if len(api.Tags) > 0 {
				module = api.Tags[0]
			}
		}

		op := &syslogrpc.VisitLogNewReq{
			PageName: module,
		}

		_, err := m.sr.AddVisitLog(r.Context(), op)
		if err != nil {
			logx.Errorf("VisitLogMiddleware Handle AddVisitLog err: %v", err)
		}
	}
}

func (m *VisitLogMiddleware) getApi(path string, method string) *spec.Operation {
	sp := m.sp
	for k, v := range sp.Paths.Paths {
		if k == path {
			switch method {
			case http.MethodGet:
				return v.Get
			case http.MethodPost:
				return v.Post
			case http.MethodPut:
				return v.Put
			case http.MethodDelete:
				return v.Delete
			case http.MethodPatch:
				return v.Patch
			case http.MethodOptions:
				return v.Options
			case http.MethodHead:
				return v.Head
			}
		}
	}

	return nil
}
