package notify_message

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/notification/notify_message"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 查询已发布但未投递的通知消息
func QueryUndeliveredNotifyMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryUndeliveredNotifyMessagesReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.Response(r, w, nil, err)
			return
		}

		l := notify_message.NewQueryUndeliveredNotifyMessagesLogic(r.Context(), svcCtx)
		resp, err := l.QueryUndeliveredNotifyMessages(&req)
		responsex.Response(r, w, resp, err)
	}
}
