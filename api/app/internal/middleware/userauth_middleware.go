package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/vkit/adapter/storex/tokenstore"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

type UserAuthMiddleware struct {
	verifier tokenstore.TokenStore
}

func NewUserAuthMiddleware(verifier tokenstore.TokenStore) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		verifier: verifier,
	}
}

func (m *UserAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("UserAuthMiddleware Handle")
		var token string
		var uid string

		uid = r.Header.Get(bizheader.HeaderUid)
		token = r.Header.Get(bizheader.HeaderToken)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderToken)))
			return
		}

		err := m.verifier.ValidateToken(uid, token)
		if err != nil {
			if errors.Is(err, tokenstore.ErrTokenExpired) {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeLoginExpired, err.Error()))
				return
			}
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, err.Error()))
			return
		}

		next(w, r)
	}
}
