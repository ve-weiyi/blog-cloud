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
		if len(respBs) > 1000 {
			respBs = []byte("response too long")
		}
		logx.WithContext(ctx).Infow("grpc server request info",
			logx.LogField{Key: "full_method", Value: info.FullMethod},
			logx.LogField{Key: "grpc_request", Value: string(reqBs)},
			logx.LogField{Key: "grpc_response", Value: string(respBs)},
			logx.LogField{Key: "err", Value: err},
		)
	}()
	return resp, err
}
