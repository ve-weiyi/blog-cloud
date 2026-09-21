package me

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/logic/account/me"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 恢复账号（冷静期内）
func ReactivateAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReactivateAccountReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := me.NewReactivateAccountLogic(r.Context(), svcCtx)
		resp, err := l.ReactivateAccount(&req)
		responsex.Response(r, w, resp, err)
	}
}
