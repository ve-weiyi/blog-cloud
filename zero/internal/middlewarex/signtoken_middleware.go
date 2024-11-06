package middlewarex

import (
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/constant"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
	"github.com/ve-weiyi/ve-blog-golang/zero/internal/responsex"
)

type SignTokenMiddleware struct {
}

func NewSignTokenMiddleware() *SignTokenMiddleware {
	return &SignTokenMiddleware{}
}

// 未登录token
// 未登录时，token = md5(tm,ts)
func (m *SignTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Infof("SignTokenMiddleware Handle")
		tk := r.Header.Get(constant.HeaderToken)
		tm := r.Header.Get(constant.HeaderTerminal)
		ts := r.Header.Get(constant.HeaderTimestamp)

		// 请求头缺少参数
		if tk == "" || tm == "" || ts == "" {
			responsex.Response(r, w, nil, apierr.NewApiError(codex.CodeUserNotPermission, "无效请求,缺少签名"))
			return
		}
		// 判断 token = md5(tm,ts)
		if tk != crypto.Md5v(tm, ts) {
			responsex.Response(r, w, nil, apierr.NewApiError(codex.CodeUserNotPermission, "无效请求,签名错误"))
			return
		}

		next.ServeHTTP(w, r)
	}
}
