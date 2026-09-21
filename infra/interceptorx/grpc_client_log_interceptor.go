package interceptorx

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func ClientLogInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)

	defer func() {
		reqBs, _ := json.Marshal(req)
		respBs, _ := json.Marshal(reply)
		if len(reqBs) > 500 {
			reqBs, _ = json.Marshal(map[string]any{
				"message": "request too large to log",
				"size":    len(reqBs),
				"body":    string(reqBs[:200]),
			})
		}
		if len(respBs) > 500 {
			respBs, _ = json.Marshal(map[string]any{
				"message": "response too large to log",
				"size":    len(respBs),
				"body":    string(respBs[:200]),
			})
		}
		logx.WithContext(ctx).Infow("grpc client log info",
			logx.LogField{Key: "full_method", Value: method},
			logx.LogField{Key: "grpc_request", Value: string(reqBs)},
			logx.LogField{Key: "grpc_response", Value: string(respBs)},
			logx.LogField{Key: "err", Value: err},
		)
	}()

	return err
}
