package interceptorx

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerErrorInterceptor transfer a error to status error, It must be the first server interceptor.
func ServerErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		log.Println("grpc server error", err)
		log.Println("grpc server status", st)
		return resp, st.Err()
	}

	return resp, nil
}
