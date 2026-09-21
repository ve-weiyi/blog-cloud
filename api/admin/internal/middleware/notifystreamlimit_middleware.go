package middleware

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/sse"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/notifyx"
	"github.com/ve-weiyi/blog-cloud/infra/responsex"
)

// DefaultStreamLimit 单实例 SSE 连接上限。后台管理端的并发连接是个位数，
// 这个值只是防止异常客户端把连接堆到失控。
const DefaultStreamLimit = 200

// NotifyStreamLimitMiddleware 限制单实例的 SSE 连接数
//
// 必须在**中间件**里拒绝：SSE handler 一旦开始写流就已接管响应头，
// 之后再也发不出信封错误，客户端只会看到一个立刻关闭的空流。
type NotifyStreamLimitMiddleware struct {
	broker *sse.Broker[notifyx.Event]
	limit  int
}

// NewNotifyStreamLimitMiddleware 创建限流中间件；limit <= 0 时取 DefaultStreamLimit
func NewNotifyStreamLimitMiddleware(broker *sse.Broker[notifyx.Event], limit int) *NotifyStreamLimitMiddleware {
	if limit <= 0 {
		limit = DefaultStreamLimit
	}
	return &NotifyStreamLimitMiddleware{broker: broker, limit: limit}
}

func (m *NotifyStreamLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if current := m.broker.Len(); current >= m.limit {
			logx.Errorf("SSE 连接数已达上限 %d（当前 %d），拒绝新连接", m.limit, current)
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeServiceUnavailable,
				fmt.Sprintf("推送连接数已达上限 %d，请稍后重试", m.limit)))
			return
		}
		next.ServeHTTP(w, r)
	}
}
