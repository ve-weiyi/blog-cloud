package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
)

// IP限流
func LimitIP(cache *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		v, err := cache.Get(context.Background(), key).Result()
		if err != nil && v != "" {
			cache.SetEx(context.Background(), key, 1, 300)
		}

		// 短时间内请求10次
		if cast.ToInt(v) > 10 {
			c.JSON(http.StatusOK, apierr.NewApiError(codex.CodeTooManyRequests, "操作频繁,请在5分钟后再试"))
			c.Abort()
			return
		}

		c.Next()
	}
}
