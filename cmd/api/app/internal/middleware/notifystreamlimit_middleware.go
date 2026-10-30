package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/http/headerx"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

const (
	// DefaultStreamGlobalLimit 单实例 SSE 连接总数上限
	DefaultStreamGlobalLimit = 500
	// DefaultStreamUserLimit 单个用户的并发连接上限，覆盖多标签页与多设备
	DefaultStreamUserLimit = 3
)

// NotifyStreamLimitMiddleware 限制读者端 SSE 连接：按用户与按实例双向限制。
//
// 必须在中间件里拒绝：SSE handler 一旦开始写流就已接管响应头，
// 之后再也发不出信封错误，客户端只会看到一个立刻关闭的空流。
type NotifyStreamLimitMiddleware struct {
	mu          sync.Mutex
	perUser     map[string]int
	active      int
	globalLimit int
	userLimit   int
}

func NewNotifyStreamLimitMiddleware(globalLimit, userLimit int) *NotifyStreamLimitMiddleware {
	if globalLimit <= 0 {
		globalLimit = DefaultStreamGlobalLimit
	}
	if userLimit <= 0 {
		userLimit = DefaultStreamUserLimit
	}
	return &NotifyStreamLimitMiddleware{
		perUser:     make(map[string]int),
		globalLimit: globalLimit,
		userLimit:   userLimit,
	}
}

func (m *NotifyStreamLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, _ := r.Context().Value(headerx.HeaderUid).(string)
		if !m.acquire(uid) {
			logx.Errorf("SSE 连接数超限，拒绝新连接: user=%s", uid)
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeServiceUnavailable,
				fmt.Sprintf("推送连接数已达上限（单用户 %d，单实例 %d），请稍后重试", m.userLimit, m.globalLimit)))
			return
		}
		defer m.release(uid)

		next.ServeHTTP(w, r)
	}
}

// acquire 先占名额再放行；任一上限命中都返回 false。
func (m *NotifyStreamLimitMiddleware) acquire(uid string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.active >= m.globalLimit {
		return false
	}
	if uid != "" && m.perUser[uid] >= m.userLimit {
		return false
	}
	m.active++
	if uid != "" {
		m.perUser[uid]++
	}
	return true
}

func (m *NotifyStreamLimitMiddleware) release(uid string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.active--
	if uid == "" {
		return
	}
	if m.perUser[uid] <= 1 {
		delete(m.perUser, uid)
		return
	}
	m.perUser[uid]--
}
