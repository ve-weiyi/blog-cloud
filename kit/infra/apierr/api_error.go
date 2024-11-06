package apierr

import "fmt"

// ApiError 是一个API错误的结构体
type ApiError struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Error 返回错误的消息
func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) Details() string {
	return fmt.Sprintf("code:%d, message:'%s'", e.Code, e.Message)
}

func (e *ApiError) WrapMessage(err string) *ApiError {
	ne := &ApiError{
		Code:    e.Code,
		Message: err,
		Data:    nil,
	}
	return ne
}

// NewApiError 创建一个新的API错误
func NewApiError(code int64, message string) *ApiError {
	return &ApiError{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
