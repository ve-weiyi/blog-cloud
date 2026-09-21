package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/middleware/visitx"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type VisitLogMiddleware struct {
	enforcer visitx.Enforcer
	sr       syslogservice.SyslogService
}

func NewVisitLogMiddleware(enforcer visitx.Enforcer, sr syslogservice.SyslogService) *VisitLogMiddleware {
	return &VisitLogMiddleware{enforcer: enforcer, sr: sr}
}

func (m *VisitLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		pageName, ok := m.enforcer.IsVisitPage(r.URL.Path, r.Method)
		if !ok {
			return
		}

		uid, _ := metax.GetApiUserIdFromCtx(r.Context())
		did, _ := metax.GetApiDeviceIdFromCtx(r.Context())

		_, err := m.sr.CreateVisitLog(r.Context(), &syslogservice.CreateVisitLogRequest{
			UserId:   uid,
			DeviceId: did,
			PageName: pageName,
		})
		if err != nil {
			logx.Errorf("VisitLog CreateVisitLog err: %v", err)
		}
	}
}
