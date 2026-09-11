package bizerr

import (
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const bizErrorDomain = "ve-blog-biz-error"

// WithBizError 将 BizError 的业务错误码编码到 gRPC status 的 ErrorInfo details 中
func WithBizError(err error) error {
	var bizErr *BizError
	if !asBizError(err, &bizErr) {
		return status.New(codes.Internal, err.Error()).Err()
	}

	st := status.New(codes.Internal, bizErr.Message)
	st, _ = st.WithDetails(&errdetails.ErrorInfo{
		Domain: bizErrorDomain,
		Reason: strconv.FormatInt(bizErr.Code, 10),
	})
	return st.Err()
}

// FromStatus 从 gRPC status 的 ErrorInfo details 中解码 BizError
// 如果不是 BizError 则返回基于 status message 的普通 error
func FromStatus(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, detail := range st.Details() {
		if info, ok := detail.(*errdetails.ErrorInfo); ok && info.Domain == bizErrorDomain {
			code, e := strconv.ParseInt(info.Reason, 10, 64)
			if e != nil {
				break
			}
			return NewBizError(code, st.Message())
		}
	}

	return fmt.Errorf(st.Message())
}

func asBizError(err error, target **BizError) bool {
	if bizErr, ok := err.(*BizError); ok {
		*target = bizErr
		return true
	}
	return false
}

// BizError 是一个业务错误的结构体
type BizError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"` // 原始错误，不序列化到JSON
}

// Error 返回错误的消息
func (e *BizError) Error() string {
	return e.Message
}

// Details 返回详细的错误信息
func (e *BizError) Details() string {
	if e.Err != nil {
		return fmt.Sprintf("code:%d, message:'%s', err:%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code:%d, message:'%s'", e.Code, e.Message)
}

// Unwrap 返回原始错误，支持 errors.Is 和 errors.As
func (e *BizError) Unwrap() error {
	return e.Err
}

// NewBizError 创建一个新的业务错误
func NewBizError(code int64, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

// WrapBizError 包装一个错误为业务错误
func WrapBizError(code int64, message string, err error) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
