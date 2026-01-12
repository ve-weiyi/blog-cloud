package _test_test

import (
	"context"
	"log"
	"testing"

	"go.opentelemetry.io/otel/sdk/trace"

	accountrpclogic "github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/logic/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
)

func TestGetUserAreasAnalysis(t *testing.T) {
	// context注入traceID
	ctx := context.Background()
	tracer := trace.NewTracerProvider().Tracer("test")
	ctx, span := tracer.Start(ctx, "trace")
	defer span.End()

	tsc := svc.NewTestServiceContext()
	in := &accountrpc.AnalysisUserReq{
		UserType: 1,
	}

	out, err := accountrpclogic.NewAnalysisUserLogic(ctx, tsc).AnalysisUser(in)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(out)
}
