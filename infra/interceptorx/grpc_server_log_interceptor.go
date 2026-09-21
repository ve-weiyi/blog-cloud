package interceptorx

import (
	"context"
	"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
)

func ServerLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	defer func() {
		reqBs, _ := json.Marshal(req)
		respBs, _ := json.Marshal(resp)
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
		logx.WithContext(ctx).Infow("grpc server log info",
			logx.LogField{Key: "full_method", Value: info.FullMethod},
			logx.LogField{Key: "grpc_request", Value: string(reqBs)},
			logx.LogField{Key: "grpc_response", Value: string(respBs)},
			logx.LogField{Key: "err", Value: err},
		)
	}()
	return resp, err
}
