package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr/codex"
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/glog"
)

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查断开的连接，因为它不是保证紧急堆栈跟踪的真正条件。
				var brokenPipe bool
				// OpError 是 net 包中的函数通常返回的错误类型。它描述了错误的操作、网络类型和地址。
				if ne, ok := err.(*net.OpError); ok {
					// SyscallError 记录来自特定系统调用的错误。
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// DumpRequest 以 HTTP/1.x 连线形式返回给定的请求
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					glog.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接死了，我们就不能给它写状态
					c.Error(err.(error))
					c.Abort() // 终止该中间件
					return
				}

				if stack {
					glog.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 返回调用它的goroutine的格式化堆栈跟踪。
					)
				} else {
					glog.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, apierr.NewApiError(codex.CodeInternalServerError, "系统错误"))
				return
			}
		}()
		c.Next()
	}
}
