package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aiquota"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/config"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/http/requestx"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
)

const (
	aiConcurrencyKey       = "blog:ai:concurrent"
	aiQuotaKeyTTL          = 48 * time.Hour
	aiConcurrencyReleaseTS = 3 * time.Second
)

// acquireConcurrency 原子地清理过期租约、检查上限并登记本轮请求。
var acquireConcurrency = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local expireAt = tonumber(ARGV[2])
local member = ARGV[3]
local limit = tonumber(ARGV[4])

redis.call('ZREMRANGEBYSCORE', key, '-inf', now)
if redis.call('ZCARD', key) >= limit then
  return 0
end
redis.call('ZADD', key, expireAt, member)
return 1
`)

// AiRateLimitMiddleware 对问答流做日配额、来源限额与全站并发控制。
//
// 身份由 OptionalUserAuth 写入 context；Redis 不可用时放行并记录，
// 配额是保护性约束，不应因为计数依赖故障把问答整个挡住。
type AiRateLimitMiddleware struct {
	rdb  *redis.Client
	conf config.AiQuotaConf
}

func NewAiRateLimitMiddleware(rdb *redis.Client, conf config.AiQuotaConf) *AiRateLimitMiddleware {
	return &AiRateLimitMiddleware{rdb: rdb, conf: conf}
}

func (m *AiRateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		identity, ok := aictx.IdentityFrom(ctx)
		if !ok {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少身份信息"))
			return
		}

		limit := aiquota.Limit(m.conf, identity.Guest)
		subject := aiquota.Subject(identity)
		day := time.Now().Format("20060102")

		// 先判并发：被并发拒绝的请求不应消耗日额度
		var release func()
		if m.conf.MaxConcurrent > 0 && m.rdb != nil {
			member, err := newAIRequestID()
			if err != nil {
				logx.WithContext(ctx).Errorf("生成 AI 请求标识失败: %v", err)
			} else {
				acquired, err := m.acquire(ctx, member)
				switch {
				case err != nil:
					logx.WithContext(ctx).Errorf("AI 并发计数失败: %v", err)
				case !acquired:
					responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "当前使用人数较多，请稍后再试"))
					return
				default:
					release = func() { m.release(member) }
				}
			}
		}
		if release != nil {
			defer release()
		}

		dailyKey := aiquota.DailyKey(subject, day)
		used, err := m.incrDaily(ctx, dailyKey)
		counted := err == nil
		if err != nil {
			logx.WithContext(ctx).Errorf("AI 配额计数失败: %v", err)
		} else if limit > 0 && used > limit {
			m.refundDaily(dailyKey)
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "今日额度已用完"))
			return
		}

		if m.conf.IPDaily > 0 {
			ip := requestx.ClientIP(r)
			if ip != "" {
				ipKey := aiquota.IPDailyKey(ip, day)
				ipUsed, err := m.incrDaily(ctx, ipKey)
				if err != nil {
					logx.WithContext(ctx).Errorf("AI 来源配额计数失败: %v", err)
				} else if ipUsed > m.conf.IPDaily {
					m.refundDaily(ipKey)
					if counted {
						m.refundDaily(dailyKey)
					}
					responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeRateLimited, "当前来源请求过于频繁"))
					return
				}
			}
		}

		left := int64(-1) // 不限额度
		if limit > 0 {
			left = limit - used
			if left < 0 {
				left = 0
			}
		}
		next(w, r.WithContext(aictx.WithQuotaLeft(ctx, left)))
	}
}

func (m *AiRateLimitMiddleware) incrDaily(ctx context.Context, key string) (int64, error) {
	if m.rdb == nil {
		return 0, nil
	}
	pipe := m.rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, aiQuotaKeyTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

func (m *AiRateLimitMiddleware) acquire(ctx context.Context, member string) (bool, error) {
	now := time.Now()
	lease := m.conf.Lease
	if lease <= 0 {
		lease = 3 * time.Minute
	}
	result, err := acquireConcurrency.Run(ctx, m.rdb,
		[]string{aiConcurrencyKey},
		now.Unix(),
		now.Add(lease).Unix(),
		member,
		m.conf.MaxConcurrent,
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

// refundDaily 撤销一次已计入的日配额；计数失败时保持原值。
func (m *AiRateLimitMiddleware) refundDaily(key string) {
	if m.rdb == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), aiConcurrencyReleaseTS)
	defer cancel()
	if err := m.rdb.Decr(ctx, key).Err(); err != nil {
		logx.Errorf("回退 AI 配额失败: %v", err)
	}
}

func (m *AiRateLimitMiddleware) release(member string) {
	ctx, cancel := context.WithTimeout(context.Background(), aiConcurrencyReleaseTS)
	defer cancel()
	if err := m.rdb.ZRem(ctx, aiConcurrencyKey, member).Err(); err != nil {
		logx.Errorf("释放 AI 并发租约失败: %v", err)
	}
}

func newAIRequestID() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("middleware: random request id: %w", err)
	}
	return hex.EncodeToString(buf), nil
}
