package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/responsex"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/tokenx"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"
)

type AdminTokenMiddleware struct {
	verifier tokenx.TokenManager
}

func NewAdminTokenMiddleware(verifier tokenx.TokenManager) *AdminTokenMiddleware {
	return &AdminTokenMiddleware{
		verifier: verifier,
	}
}

func (m *AdminTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("AdminTokenMiddleware Handle")
		var token string
		var uid string

		token = r.Header.Get(bizheader.HeaderAuthorization)
		uid = r.Header.Get(bizheader.HeaderUid)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderAuthorization)))
			return
		}

		err := m.verifier.ValidateToken(uid, token)
		if err != nil {
			if errors.Is(err, tokenx.ErrTokenExpired) {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUserLoginExpired, err.Error()))
				return
			}
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUserUnLogin, err.Error()))
			return
		}

		next.ServeHTTP(w, r)
	}
}
