package httperr

import "fmt"

// HttpError 是一个http错误的结构体
type HttpError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Error 返回错误的消息
func (e *HttpError) Error() string {
	return e.Message
}

func (e *HttpError) Details() string {
	return fmt.Sprintf("code:%d, message:'%s'", e.Code, e.Message)
}

func (e *HttpError) WrapMessage(err string) *HttpError {
	ne := &HttpError{
		Code:    e.Code,
		Message: err,
	}
	return ne
}

// NewHttpError 创建一个新的API错误
func NewHttpError(code int64, message string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}
