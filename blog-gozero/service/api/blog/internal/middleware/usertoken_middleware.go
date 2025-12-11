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

type UserTokenMiddleware struct {
	verifier tokenx.TokenHolder
}

func NewUserTokenMiddleware(verifier tokenx.TokenHolder) *UserTokenMiddleware {
	return &UserTokenMiddleware{
		verifier: verifier,
	}
}

func (m *UserTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("SignTokenMiddleware Handle")
		var token string
		var uid string

		uid = r.Header.Get(restx.HeaderUid)
		token = r.Header.Get(restx.HeaderToken)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderToken)))
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
