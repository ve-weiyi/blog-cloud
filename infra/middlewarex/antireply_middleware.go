package middlewarex

import (
	"net/http"
	"time"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

type AntiReplyMiddleware struct {
}

func NewAntiReplyMiddleware() *AntiReplyMiddleware {
	return &AntiReplyMiddleware{}
}

// 防重放
func (m *AntiReplyMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Debugf("AntiReplyMiddleware Handle")
		ts := r.Header.Get(bizheader.HeaderTimestamp)
		if ts != "" {
			now := time.Now().Unix()
			if now-cast.ToInt64(ts) > 3600 {
				responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRequestSignInvalid, "timestamp is expired"))
				return
			}
		}

		next.ServeHTTP(w, r)
	}
}
