package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/account/user"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 重置用户密码
func ResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetUserPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := user.NewResetUserPasswordLogic(r.Context(), svcCtx)
		resp, err := l.ResetUserPassword(&req)
		responsex.Response(r, w, resp, err)
	}
}
