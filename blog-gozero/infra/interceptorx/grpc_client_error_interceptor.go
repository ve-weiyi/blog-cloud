package interceptorx

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ClientErrorInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if ok {
		details := st.Details()
		if len(details) > 0 {
			return err
		}
	}

	return err
}
