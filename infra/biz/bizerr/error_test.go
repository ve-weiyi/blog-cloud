package bizerr

import (
	"errors"
	"fmt"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
)

// TestCanonicalByCodeCoversAllCodes 守住"新增标识必须同时确定它的 gRPC canonical 对应"
// 这条约束：映射表与 bizcode.All 任一方向漏了都会在这里失败。
func TestCanonicalByCodeCoversAllCodes(t *testing.T) {
	declared := make(map[string]bool, len(bizcode.All))
	for _, code := range bizcode.All {
		if declared[code] {
			t.Errorf("bizcode.All 存在重复标识: %s", code)
		}
		declared[code] = true
	}

	for _, code := range bizcode.All {
		if _, ok := canonicalByCode[code]; !ok {
			t.Errorf("标识 %s 缺少 canonical 映射，请补入 bizerr.canonicalByCode", code)
		}
	}

	for code := range canonicalByCode {
		if !declared[code] {
			t.Errorf("映射表存在未声明的标识 %s，请补入 bizcode.All 或删除该映射", code)
		}
	}
}

// TestWithBizErrorCarriesCanonical 守住 gRPC 一跳的 canonical 契约：
// canonical code 由标识确定，不能一概压成 Internal。
func TestWithBizErrorCarriesCanonical(t *testing.T) {
	cases := []struct {
		bizCode string
		want    codes.Code
	}{
		{bizcode.CodeResourceNotFound, codes.NotFound},
		{bizcode.CodeResourceAlreadyExist, codes.AlreadyExists},
		{bizcode.CodeResourceStatusNotAllowed, codes.FailedPrecondition},
		{bizcode.CodePreconditionFailed, codes.Aborted},
		{bizcode.CodeRateLimited, codes.ResourceExhausted},
		{bizcode.CodeServiceUnavailable, codes.Unavailable},
		{bizcode.CodeServiceTimeout, codes.DeadlineExceeded},
		{bizcode.CodeUnauthenticated, codes.Unauthenticated},
		{bizcode.CodeNoPermission, codes.PermissionDenied},
		{bizcode.CodeParamFormat, codes.InvalidArgument},
	}

	for _, c := range cases {
		t.Run(c.bizCode, func(t *testing.T) {
			err := WithBizError(NewBizError(c.bizCode, "boom"))

			if got := status.Code(err); got != c.want {
				t.Errorf("canonical code = %v, want %v", got, c.want)
			}
		})
	}
}

// TestWithBizErrorCarriesViolations 守住字段级明细在 gRPC 侧的承载：
// Violations 要挂到 details[].BadRequest，而不是只留在 message 里。
func TestWithBizErrorCarriesViolations(t *testing.T) {
	bizErr := NewBizError(bizcode.CodeParamFormat, "参数校验不通过")
	bizErr.Violations = []FieldViolation{
		{Field: "title", Description: "不能为空"},
		{Field: "tags[2].name", Description: "长度超过 32"},
	}

	st := status.Convert(WithBizError(bizErr))

	var badRequest *errdetails.BadRequest
	for _, detail := range st.Details() {
		if br, ok := detail.(*errdetails.BadRequest); ok {
			badRequest = br
		}
	}
	if badRequest == nil {
		t.Fatal("details 里没有 BadRequest")
	}
	if len(badRequest.FieldViolations) != 2 {
		t.Fatalf("FieldViolations 条数 = %d, want 2", len(badRequest.FieldViolations))
	}
	if badRequest.FieldViolations[0].Field != "title" {
		t.Errorf("FieldViolations[0].Field = %q, want %q", badRequest.FieldViolations[0].Field, "title")
	}
}

// TestWithBizErrorUnwrapsWrappedError 守住包装过的 BizError 也要被识别：
// fmt.Errorf("%w") 包装后必须仍能取到 canonical，而不是被压成 Internal。
func TestWithBizErrorUnwrapsWrappedError(t *testing.T) {
	wrapped := fmt.Errorf("调用下游失败: %w", NewBizError(bizcode.CodeResourceNotFound, "文章不存在"))

	if got := status.Code(WithBizError(wrapped)); got != codes.NotFound {
		t.Errorf("包装后的 BizError 未被识别，canonical = %v, want %v", got, codes.NotFound)
	}
}

// TestFromStatusFallsBackByCanonical 守住兜底：下游只给了 canonical code 时，
// 网关仍要拿到一个标识——否则客户端拿不到 code。
func TestFromStatusFallsBackByCanonical(t *testing.T) {
	cases := []struct {
		canonical codes.Code
		want      string
	}{
		{codes.Unauthenticated, bizcode.CodeUnauthenticated},
		{codes.PermissionDenied, bizcode.CodeNoPermission},
		{codes.InvalidArgument, bizcode.CodeInvalidParam},
		{codes.OutOfRange, bizcode.CodeInvalidParam},
		{codes.NotFound, bizcode.CodeResourceNotFound},
		{codes.AlreadyExists, bizcode.CodeResourceAlreadyExist},
		{codes.FailedPrecondition, bizcode.CodeResourceStatusNotAllowed},
		{codes.Aborted, bizcode.CodePreconditionFailed},
		{codes.ResourceExhausted, bizcode.CodeRateLimited},
		{codes.Unavailable, bizcode.CodeServiceUnavailable},
		{codes.DeadlineExceeded, bizcode.CodeServiceTimeout},
		{codes.Unknown, bizcode.CodeInternalServerError},
		{codes.DataLoss, bizcode.CodeInternalServerError},
		{codes.Canceled, bizcode.CodeInternalServerError},
	}

	for _, c := range cases {
		t.Run(c.canonical.String(), func(t *testing.T) {
			got := FromStatus(status.Error(c.canonical, "downstream failed"))

			var bizErr *BizError
			if !errors.As(got, &bizErr) {
				t.Fatalf("FromStatus 未返回 *BizError，得到 %T", got)
			}
			if bizErr.Code != c.want {
				t.Errorf("兜底标识 = %s, want %s", bizErr.Code, c.want)
			}
		})
	}
}

// TestFromStatusNil 守住 nil 透传：空错误不得被兜底成一个错误。
func TestFromStatusNil(t *testing.T) {
	if got := FromStatus(nil); got != nil {
		t.Errorf("FromStatus(nil) = %v, want nil", got)
	}
}

// TestFromStatusPassesThroughNonStatusError 守住边界：非 gRPC status 的错误原样返回，
// 不被包装成业务错误——它不是下游的语义，不该被赋予标识。
func TestFromStatusPassesThroughNonStatusError(t *testing.T) {
	origin := errors.New("connection refused")

	if got := FromStatus(origin); got != origin {
		t.Errorf("非 status 错误应原样返回，得到 %v", got)
	}
}

// TestFromStatusPrefersErrorInfo 守住优先级：下游已给出标识时原样带出，不走 canonical 兜底。
func TestFromStatusPrefersErrorInfo(t *testing.T) {
	st := status.New(codes.Internal, "boom")
	st, err := st.WithDetails(&errdetails.ErrorInfo{
		Domain: bizErrorDomain,
		Reason: bizcode.CodeResourceNotFound,
	})
	if err != nil {
		t.Fatalf("构造 status 失败: %v", err)
	}

	got := FromStatus(st.Err())

	var bizErr *BizError
	if !errors.As(got, &bizErr) {
		t.Fatalf("FromStatus 未返回 *BizError，得到 %T", got)
	}
	if bizErr.Code != bizcode.CodeResourceNotFound {
		t.Errorf("标识 = %s, want %s（应优先用 ErrorInfo.reason 而非 canonical 兜底）", bizErr.Code, bizcode.CodeResourceNotFound)
	}
}
