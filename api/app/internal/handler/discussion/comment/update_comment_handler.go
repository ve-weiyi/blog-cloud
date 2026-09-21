package comment

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/logic/discussion/comment"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 更新评论
func UpdateCommentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateCommentReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := comment.NewUpdateCommentLogic(r.Context(), svcCtx)
		resp, err := l.UpdateComment(&req)
		responsex.Response(r, w, resp, err)
	}
}
