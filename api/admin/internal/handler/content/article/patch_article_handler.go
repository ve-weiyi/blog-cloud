package article

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/logic/content/article"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 部分更新文章（置顶状态 / 删除状态）
func PatchArticleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PatchArticleReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := article.NewPatchArticleLogic(r.Context(), svcCtx)
		resp, err := l.PatchArticle(&req)
		responsex.Response(r, w, resp, err)
	}
}
