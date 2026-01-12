package interceptorx

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerErrorInterceptor transfer a error to status error, It must be the first server interceptor.
func ServerErrorInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	resp, err := handler(ctx, req)
	if err != nil {
		st := status.New(codes.Internal, err.Error())
		logx.WithContext(ctx).Errorf("grpc server error info, full_method=%s, err=%+v", info.FullMethod, err)
		return resp, st.Err()
	}

	return resp, nil
}
