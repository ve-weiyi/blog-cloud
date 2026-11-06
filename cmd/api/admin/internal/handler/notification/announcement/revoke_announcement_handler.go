package announcement

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/logic/notification/announcement"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

// 撤回公告
func RevokeAnnouncementHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RevokeAnnouncementReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := announcement.NewRevokeAnnouncementLogic(r.Context(), svcCtx)
		resp, err := l.RevokeAnnouncement(&req)
		responsex.Response(r, w, resp, err)
	}
}
