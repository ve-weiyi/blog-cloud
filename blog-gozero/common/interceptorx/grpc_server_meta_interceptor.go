package interceptorx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 输出请求的元数据，方便其他rpc服务获取
func ServerMetaInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		var pairs []string
		for k, v := range md {
			//logx.Infof("ServerMetaInterceptor get k=%s, v=%+v", k, v)
			for _, value := range v {
				pairs = append(pairs, k, value)
			}
		}
		ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
	}

	return handler(ctx, req)
}
