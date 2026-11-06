package quota

import (
	"net/http"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/logic/ai/quota"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

// 查询今日剩余额度
func GetAiQuotaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := quota.NewGetAiQuotaLogic(r.Context(), svcCtx)
		resp, err := l.GetAiQuota()
		responsex.Response(r, w, resp, err)
	}
}
