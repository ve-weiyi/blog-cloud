package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/sdk/trace"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := c.Request.Context()

		// context注入traceID
		tracer := trace.NewTracerProvider().Tracer(ip)
		ctx, span := tracer.Start(ctx, "trace")
		defer span.End()

		// 将 Trace ID 存入 context
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
