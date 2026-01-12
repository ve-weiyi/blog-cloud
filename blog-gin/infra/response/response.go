package response

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/trace"

	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizcode"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
)

type Body struct {
	Code    int64       `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id,omitempty"`
}

func Response(c *gin.Context, code int64, msg string, data interface{}) {
	traceID := trace.SpanContextFromContext(c.Request.Context())

	obj := Body{
		Code:    code,
		Msg:     msg,
		Data:    data,
		TraceId: traceID.TraceID().String(),
	}
	c.JSON(bizcode.CodeSuccess, obj)
}

func ResponseOk(c *gin.Context, data interface{}) {
	Response(c, bizcode.CodeSuccess, "Operation successful", data)
}

func ResponseError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// 1. 返回错误信息多语言
	// 2. 忽略详细的内部错误信息
	// 3. 处理不同的业务逻辑
	var bizErr *bizerr.BizError
	var mysqlErr *mysql.MySQLError

	switch {
	case errors.As(err, &bizErr):
		Response(c, bizErr.Code, bizErr.Error(), nil)
	case errors.As(err, &mysqlErr):
		Response(c, bizcode.CodeSqlQueryError, SqlErrorI18n(mysqlErr), nil)
	default:
		Response(c, bizcode.CodeInternalServerError, err.Error(), nil)
	}
}
