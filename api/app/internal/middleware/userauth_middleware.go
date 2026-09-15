package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
	"github.com/ve-weiyi/blog-cloud/infra/tokenx"
)

type UserAuthMiddleware struct {
	verifier tokenx.Manager
}

func NewUserAuthMiddleware(verifier tokenx.Manager) *UserAuthMiddleware {
	return &UserAuthMiddleware{
		verifier: verifier,
	}
}

func (m *UserAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("UserAuthMiddleware Handle")
		var token string
		var uid string
		var deviceId string

		uid = r.Header.Get(bizheader.HeaderUid)
		token = r.Header.Get(bizheader.HeaderToken)
		deviceId = r.Header.Get(bizheader.HeaderXDeviceId)

		// 请求头缺少参数
		if uid == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderUid)))
			return
		}

		if token == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderToken)))
			return
		}

		err := m.verifier.Validate(r.Context(), uid, deviceId, token)
		if err != nil {
			if errors.Is(err, tokenx.ErrTokenExpired) {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeLoginExpired, err.Error()))
				return
			}
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, err.Error()))
			return
		}

		next(w, r)
	}
}
