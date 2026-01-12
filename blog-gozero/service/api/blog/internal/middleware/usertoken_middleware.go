package middleware

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/responsex"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"
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

		uid = r.Header.Get(bizheader.HeaderUid)
		token = r.Header.Get(bizheader.HeaderToken)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderToken)))
			return
		}

		err := m.verifier.VerifyToken(r.Context(), token, uid)
		if err != nil {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUserLoginExpired, err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	}
}
