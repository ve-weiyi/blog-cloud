package middlewarex

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// RateLimitMiddleware 接口频率限制中间件。
type RateLimitMiddleware struct {
	limiter limitx.Limiter
}

func NewRateLimitMiddleware(limiter limitx.Limiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: limiter}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := fmt.Sprintf("%s:%s", httpx.GetRemoteAddr(r), r.URL.Path)
		code, err := m.limiter.Take(r.Context(), key)
		if err != nil {
			logx.Errorf("rate limit check failed: %v", err)
			next(w, r)
			return
		}
		if code == limitx.OverQuota {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "请求过于频繁"))
			return
		}
		next(w, r)
	}
}
