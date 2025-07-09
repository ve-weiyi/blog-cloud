package bizerr

import "fmt"

// 是一个业务错误的结构体
type BizError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// 返回错误的消息
func (e *BizError) Error() string {
	return e.Message
}

func (e *BizError) Details() string {
	return fmt.Sprintf("code:%d, message:'%s'", e.Code, e.Message)
}

// NewBizError 创建一个新的API错误
func NewBizError(code int64, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}
