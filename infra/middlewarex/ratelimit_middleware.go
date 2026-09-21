package middlewarex

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/middlewarex/limitx"
	"github.com/ve-weiyi/blog-cloud/infra/requestx"
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
		key := fmt.Sprintf("%s:%s", requestx.ClientIP(r), r.URL.Path)
		result, err := m.limiter.Take(r.Context(), key)
		if err != nil {
			// 计数不可得时放行：限流是保护性约束而非正确性约束，
			// 不能因为 Redis 故障就把全站请求挡在 429 后面
			logx.Errorf("rate limit check failed: %v", err)
			next(w, r)
			return
		}

		setRateLimitHeaders(w, result)

		if result.State == limitx.OverQuota {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "请求过于频繁"))
			return
		}
		next(w, r)
	}
}

// setRateLimitHeaders 回传余量口径，使客户端能提前降速而不是撞到 429 才知道；
// 超限时同时给出精确的 Retry-After。时长向上取整到秒——客户端按秒退避。
func setRateLimitHeaders(w http.ResponseWriter, result limitx.Result) {
	resetSeconds := strconv.Itoa(secondsCeil(result.Reset))

	w.Header().Set(responsex.HeaderRateLimitLimit, strconv.Itoa(result.Limit))
	w.Header().Set(responsex.HeaderRateLimitRemaining, strconv.Itoa(result.Remaining))
	w.Header().Set(responsex.HeaderRateLimitReset, resetSeconds)

	if result.State == limitx.OverQuota {
		w.Header().Set(responsex.HeaderRetryAfter, resetSeconds)
	}
}

func secondsCeil(d time.Duration) int {
	seconds := int(math.Ceil(d.Seconds()))
	if seconds < 1 {
		return 1
	}

	return seconds
}
