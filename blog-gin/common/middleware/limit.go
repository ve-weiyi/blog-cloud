package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/response"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/biz/bizerr"
)

// 限频
func Limit(rdb *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		key := c.ClientIP()
		ctx := c.Request.Context()

		// 短时间内请求10次
		if rateLimiter(rdb, ctx, key) {
			response.ResponseError(c, bizerr.NewBizError(bizerr.CodeTooManyRequests, "操作频繁,请在5分钟后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}

const (
	limit    = 10
	interval = 5 * time.Minute
)

// 使用 SETNX 和 INCR 实现限流
func rateLimiter(rdb *redis.Client, ctx context.Context, key string) bool {
	// 尝试设置过期时间
	set, err := rdb.SetNX(ctx, key, 1, interval).Result()
	if err != nil {
		return false
	}

	// 如果 key 不存在，说明这是第一个请求，设置成功，允许通过
	if set {
		return true
	}

	// 如果 key 已存在，说明在当前时间窗口内，再次请求递增计数器
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	// 如果计数器超出限制，拒绝请求
	return count <= int64(limit)
}
