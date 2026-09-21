package notify_record

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/notification/notify_record"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 批量删除投递记录
func BatchDeleteNotifyRecordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteNotifyRecordReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := notify_record.NewBatchDeleteNotifyRecordLogic(r.Context(), svcCtx)
		resp, err := l.BatchDeleteNotifyRecord(&req)
		responsex.Response(r, w, resp, err)
	}
}
