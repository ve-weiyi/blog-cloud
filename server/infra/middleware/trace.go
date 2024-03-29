package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const TraceIDKey = "X-Trace-ID"

// 链路追踪 https://coder55.com/article/159071
type spanContext struct {
	traceId string // TraceID 表示tracer的全局唯一ID
	spanId  string // SpanId 标示单个trace中某一个span的唯一ID，在trace中唯一
}

// SpanContext 保存了分布式追踪的上下文信息，包括 Trace id，Span id 以及其它需要传递到下游的内容。
type SpanContext interface {
	TraceId() string                     // get TraceId
	SpanId() string                      // get SpanId
	Visit(fn func(key, val string) bool) // 自定义操作TraceId，SpanId
}

// 一个 REST 调用或者数据库操作等，都可以作为一个 span 。 span 是分布式追踪的最小跟踪单位，一个 Trace 由多段 Span 组成。
type Span struct {
	ctx           spanContext // 传递的上下文
	serviceName   string      // 服务名
	operationName string      // 操作
	startTime     time.Time   // 开始时间戳
	flag          string      // 标记开启trace是 server 还是 client
	children      int         // 本 span fork出来的 childsnums
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头中获取 Trace ID
		traceID := c.GetHeader(TraceIDKey)

		if traceID == "" {
			// 如果请求头中不存在 Trace ID，则生成一个新的 Trace ID
			traceID = GenerateTraceID()
		}
		// 将 Trace ID 存入 context
		ctx := ContextWithFields(c.Request.Context(), zap.String(TraceIDKey, traceID))
		ctx = context.WithValue(ctx, TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		// 将 Trace ID 设置到响应头中，以便后续服务使用
		c.Header(TraceIDKey, traceID)

		c.Next()
	}
}

func GenerateTraceID() string {
	return uuid.NewString()
}

const (
	fieldsContextKey = "contextKey"
)

func ContextWithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if val := ctx.Value(fieldsContextKey); val != nil {
		if arr, ok := val.([]zap.Field); ok {
			allFields := make([]zap.Field, 0, len(arr)+len(fields))
			allFields = append(allFields, arr...)
			allFields = append(allFields, fields...)
			return context.WithValue(ctx, fieldsContextKey, allFields)
		}
	}

	return context.WithValue(ctx, fieldsContextKey, fields)
}
