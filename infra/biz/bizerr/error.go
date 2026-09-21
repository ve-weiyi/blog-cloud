package bizerr

import (
	"errors"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
)

const bizErrorDomain = "blog-biz-error"

// canonicalByCode 是业务错误标识到 gRPC canonical code 的对照锚点。
//
// canonical code 是标识在 gRPC 侧的粗粒度表达：它不携带标识的细分语义，但让标准 gRPC
// 中间件与下游能按类处理，也是 details 丢失时网关兜底出标识的依据（见 fallbackByCanonical）。
// 新增标识时必须同时登记此处，由同包测试保证与 bizcode.All 双向一致。
var canonicalByCode = map[string]codes.Code{
	bizcode.CodeSuccess: codes.OK,

	bizcode.CodeRateLimited:        codes.ResourceExhausted,
	bizcode.CodeRequestSignInvalid: codes.Unauthenticated,

	bizcode.CodeInvalidParam:      codes.InvalidArgument,
	bizcode.CodeParamMissing:      codes.InvalidArgument,
	bizcode.CodeParamFormat:       codes.InvalidArgument,
	bizcode.CodeParamValueInvalid: codes.InvalidArgument,
	bizcode.CodeVerifyCodeError:   codes.InvalidArgument,

	bizcode.CodeUnauthenticated:    codes.Unauthenticated,
	bizcode.CodeLoginExpired:       codes.Unauthenticated,
	bizcode.CodeCredentialsInvalid: codes.Unauthenticated,
	bizcode.CodeAccountDisabled:    codes.PermissionDenied,

	bizcode.CodeNoPermission:        codes.PermissionDenied,
	bizcode.CodeRoleNotMatch:        codes.PermissionDenied,
	bizcode.CodeOperationNotAllowed: codes.PermissionDenied,

	bizcode.CodeResourceNotFound:         codes.NotFound,
	bizcode.CodeResourceAlreadyExist:     codes.AlreadyExists,
	bizcode.CodeResourceStatusNotAllowed: codes.FailedPrecondition,
	bizcode.CodePreconditionFailed:       codes.Aborted,

	bizcode.CodeInternalServerError: codes.Internal,
	bizcode.CodeDatabaseError:       codes.Internal,
	// 外部服务返回无效响应在 canonical 集合里没有对应项：挑一个近似的码就是错配，
	// 故用 UNKNOWN 表示"无法归类"。
	bizcode.CodeExternalServiceError: codes.Unknown,
	bizcode.CodeServiceUnavailable:   codes.Unavailable,
	bizcode.CodeServiceTimeout:       codes.DeadlineExceeded,
}

// fallbackByCanonical 是"下游没给出标识"时的兜底表：canonical code 到最粗的标识。
//
// **本表是单向的**——只用于 details 里没有 ErrorInfo 的情形，不构成从 canonical
// 反推标识的通用依据：多个 canonical code 合并到同一标识是刻意的（它们的细分语义在
// 合并中丢失，而下游本来也没给出更细的信息）。
var fallbackByCanonical = map[codes.Code]string{
	codes.Unauthenticated:    bizcode.CodeUnauthenticated,
	codes.PermissionDenied:   bizcode.CodeNoPermission,
	codes.InvalidArgument:    bizcode.CodeInvalidParam,
	codes.OutOfRange:         bizcode.CodeInvalidParam,
	codes.NotFound:           bizcode.CodeResourceNotFound,
	codes.AlreadyExists:      bizcode.CodeResourceAlreadyExist,
	codes.FailedPrecondition: bizcode.CodeResourceStatusNotAllowed,
	codes.Aborted:            bizcode.CodePreconditionFailed,
	codes.ResourceExhausted:  bizcode.CodeRateLimited,
	codes.Unavailable:        bizcode.CodeServiceUnavailable,
	codes.DeadlineExceeded:   bizcode.CodeServiceTimeout,
	codes.Unimplemented:      bizcode.CodeInternalServerError,
	codes.Internal:           bizcode.CodeInternalServerError,
	codes.Unknown:            bizcode.CodeInternalServerError,
	codes.DataLoss:           bizcode.CodeInternalServerError,
	codes.Canceled:           bizcode.CodeInternalServerError,
}

// canonicalOf 返回标识对应的 canonical code；表外标识按 Internal 兜底。
func canonicalOf(bizCode string) codes.Code {
	if code, ok := canonicalByCode[bizCode]; ok {
		return code
	}

	return codes.Internal
}

// WithBizError 将 BizError 编码为 gRPC status：canonical code 按标识确定，
// 标识挂 details[].ErrorInfo，字段级明细挂 details[].BadRequest。
func WithBizError(err error) error {
	var bizErr *BizError
	if !errors.As(err, &bizErr) {
		return status.New(codes.Internal, err.Error()).Err()
	}

	canonical := canonicalOf(bizErr.Code)
	errorInfo := &errdetails.ErrorInfo{
		Domain: bizErrorDomain,
		Reason: bizErr.Code,
	}

	st := status.New(canonical, bizErr.Message)
	var withDetailsErr error
	if len(bizErr.Violations) > 0 {
		st, withDetailsErr = st.WithDetails(errorInfo, toBadRequest(bizErr.Violations))
	} else {
		st, withDetailsErr = st.WithDetails(errorInfo)
	}
	if withDetailsErr != nil {
		// 明细挂载失败时退化为不带 details 的 status：canonical code 仍然正确，
		// 网关侧可据此兜底出粗粒度标识，不至于把错误本身丢掉。
		return status.New(canonical, bizErr.Message).Err()
	}

	return st.Err()
}

// FromStatus 从 gRPC status 解码 BizError：优先取 details[].ErrorInfo.reason；
// 下游没给出标识时按 canonical code 兜底成最粗的标识，使调用方始终能拿到一个标识。
// 非 status 错误原样返回。
func FromStatus(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok && info.Domain == bizErrorDomain {
			return NewBizError(info.Reason, st.Message())
		}
	}

	if code, ok := fallbackByCanonical[st.Code()]; ok {
		return NewBizError(code, st.Message())
	}

	return NewBizError(bizcode.CodeInternalServerError, st.Message())
}

// toBadRequest 把字段级明细转成 gRPC 侧的 BadRequest。
func toBadRequest(violations []FieldViolation) *errdetails.BadRequest {
	items := make([]*errdetails.BadRequest_FieldViolation, 0, len(violations))
	for _, v := range violations {
		items = append(items, &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Description,
		})
	}

	return &errdetails.BadRequest{FieldViolations: items}
}

// FieldViolation 是字段级校验明细，形状对应 AIP-193 的 BadRequest.FieldViolation。
// `field` 用请求体中的字段路径，数组元素带下标。
type FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

// BizError 是一个业务错误的结构体
type BizError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // 原始错误，不序列化到JSON

	// Violations 是字段级校验明细，仅在参数校验类失败上使用。
	// 由响应出口写入信封的 data.field_violations，gRPC 侧写 details[].BadRequest。
	Violations []FieldViolation `json:"-"`
}

// Error 返回错误的消息
func (e *BizError) Error() string {
	return e.Message
}

// Details 返回详细的错误信息
func (e *BizError) Details() string {
	if e.Err != nil {
		return fmt.Sprintf("code:%s, message:'%s', err:%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code:%s, message:'%s'", e.Code, e.Message)
}

// Unwrap 返回原始错误，支持 errors.Is 和 errors.As
func (e *BizError) Unwrap() error {
	return e.Err
}

// NewBizError 创建一个新的业务错误
func NewBizError(code string, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// WrapBizError 包装一个错误为业务错误
func WrapBizError(code string, message string, err error) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
