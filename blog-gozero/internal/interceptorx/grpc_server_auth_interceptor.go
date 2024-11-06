package interceptorx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func ServerAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		var pairs []string
		for k, v := range md {
			for _, value := range v {
				pairs = append(pairs, k, value)
			}
		}
		ctx = metadata.AppendToOutgoingContext(ctx, pairs...)
	}

	return handler(ctx, req)
}
