package apierr

import "fmt"

// ApiError 定义了API错误的接口
type ApiError interface {
	Code() int
	Error() string
	Wrap(err error) *apiError
}

// apiError 是ApiError接口的一个实现
type apiError struct {
	code    int
	message string
}

// Code 返回错误的状态码
func (e *apiError) Code() int {
	return e.code
}

// Error 返回错误的消息
func (e *apiError) Error() string {
	return e.message
}

func (e *apiError) Message() string {
	return e.message
}

func (e *apiError) Wrap(err error) *apiError {
	ne := &apiError{message: fmt.Sprintf("%v,err:%v", e.message, err), code: e.code}
	return ne
}

// NewApiError 创建一个新的API错误
func NewApiError(code int, message string) ApiError {
	return &apiError{
		code:    code,
		message: message,
	}
}
