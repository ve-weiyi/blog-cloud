package home

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/logic/site/home"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 获取博客前台首页信息
func GetBlogHomeInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetBlogHomeInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := home.NewGetBlogHomeInfoLogic(r.Context(), svcCtx)
		resp, err := l.GetBlogHomeInfo(&req)
		responsex.Response(r, w, resp, err)
	}
}
