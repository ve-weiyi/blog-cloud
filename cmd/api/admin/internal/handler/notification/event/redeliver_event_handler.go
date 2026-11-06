package event

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/logic/notification/event"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

// 重新扇出事件
func RedeliverEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RedeliverEventReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := event.NewRedeliverEventLogic(r.Context(), svcCtx)
		resp, err := l.RedeliverEvent(&req)
		responsex.Response(r, w, resp, err)
	}
}
