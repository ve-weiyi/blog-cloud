package middlewarex

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/infra/responsex"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

type TerminalTokenMiddleware struct {
}

func NewTerminalTokenMiddleware() *TerminalTokenMiddleware {
	return &TerminalTokenMiddleware{}
}

func (m *TerminalTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("TerminalTokenMiddleware Handle")
		ts := r.Header.Get(bizheader.HeaderTimestamp)
		tm := r.Header.Get(bizheader.HeaderXTerminalId)
		tk := r.Header.Get(bizheader.HeaderXTerminalToken)

		// 没有客户端id
		if tm == "" {
			// 拦截
			//responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXTerminalId)))
			// 默认放行
			next.ServeHTTP(w, r)
			return
		}

		// 请求头缺少参数
		if ts == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderTimestamp)))
			return
		}

		if tk == "" {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeInvalidParam, fmt.Sprintf("request header field '%v' is missing", bizheader.HeaderXTerminalToken)))
			return
		}

		// 判断 token = md5(tm,ts)
		if tk != cryptox.Md5v(tm, ts) {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUserLoginExpired, "无效请求,游客签名错误"))
			return
		}

		next.ServeHTTP(w, r)
	}
}
