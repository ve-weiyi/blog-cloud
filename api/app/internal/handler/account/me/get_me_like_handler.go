package me

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/logic/account/me"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// 获取用户点赞集合
func GetMeLikeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMeLikeReq
		if err := httpx.Parse(r, &req); err != nil {
			responsex.RequestInvalid(r, w, err)
			return
		}

		l := me.NewGetMeLikeLogic(r.Context(), svcCtx)
		resp, err := l.GetMeLike(&req)
		responsex.Response(r, w, resp, err)
	}
}
