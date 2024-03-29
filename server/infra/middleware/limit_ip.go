package middleware

import (
	"github.com/gin-gonic/gin"
)

func LimitIP() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 如果需要限制ip,可以在这里进行限制

		c.Set("ip_address", ip)
		c.Next()
	}
}
