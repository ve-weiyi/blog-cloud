package middleware

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/common/responsex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"
)

type TerminalTokenMiddleware struct {
}

func NewTerminalTokenMiddleware() *TerminalTokenMiddleware {
	return &TerminalTokenMiddleware{}
}

func (m *TerminalTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("TerminalTokenMiddleware Handle")
		ts := r.Header.Get(restx.HeaderTimestamp)
		tm := r.Header.Get(restx.HeaderXTerminalId)
		tk := r.Header.Get(restx.HeaderXTerminalToken)

		// 请求头缺少参数
		if ts == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderTimestamp)))
			return
		}

		if tm == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderXTerminalId)))
			return
		}

		if tk == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", restx.HeaderXTerminalToken)))
			return
		}

		// 判断 token = md5(tm,ts)
		if tk != crypto.Md5v(tm, ts) {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizerr.CodeUserLoginExpired, "无效请求,游客签名错误"))
			return
		}

		next.ServeHTTP(w, r)
	}
}
