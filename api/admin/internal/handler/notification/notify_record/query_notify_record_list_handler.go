package notify_record

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/notification/notify_record"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 获取统一投递记录列表
func QueryNotifyRecordListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryNotifyRecordListReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := notify_record.NewQueryNotifyRecordListLogic(r.Context(), svcCtx)
		resp, err := l.QueryNotifyRecordList(&req)
		responsex.Response(r, w, resp, err)
	}
}
