package middleware

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/responsex"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
)

type AdminTokenMiddleware struct {
	verifier tokenx.TokenHolder
}

func NewAdminTokenMiddleware(verifier tokenx.TokenHolder) *AdminTokenMiddleware {
	return &AdminTokenMiddleware{
		verifier: verifier,
	}
}

func (m *AdminTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("JwtTokenMiddleware Handle")
		var token string
		var uid string

		token = r.Header.Get(restx.HeaderAuthorization)
		uid = r.Header.Get(restx.HeaderUid)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderAuthorization)))
			return
		}

		err := m.verifier.VerifyToken(r.Context(), token, uid)
		if err != nil {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeUserLoginExpired, err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	}
}
