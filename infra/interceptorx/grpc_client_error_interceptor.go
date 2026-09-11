package interceptorx

import (
	"context"

	"google.golang.org/grpc"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

func ClientErrorInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	return bizerr.FromStatus(err)
}
