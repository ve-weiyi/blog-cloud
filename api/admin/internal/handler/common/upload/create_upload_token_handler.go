package upload

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/common/upload"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 获取上传凭证（前端直传）
func CreateUploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateUploadTokenReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := upload.NewCreateUploadTokenLogic(r.Context(), svcCtx)
		resp, err := l.CreateUploadToken(&req)
		responsex.Response(r, w, resp, err)
	}
}
