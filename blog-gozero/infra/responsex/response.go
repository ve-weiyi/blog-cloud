package responsex

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go.opentelemetry.io/otel/trace"

	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizerr"
)

type Body struct {
	Code        int64       `json:"code"`
	Msg         string      `json:"msg"`
	Data        interface{} `json:"data,omitempty"`
	EncryptData string      `json:"encrypt_data,omitempty"`
	TraceId     string      `json:"trace_id"`
}

// Response 统一封装成功响应值.
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		var bizErr *bizerr.BizError
		var unmarshalErr *json.UnmarshalTypeError
		var mysqlErr *mysql.MySQLError

		switch {
		case errors.As(err, &bizErr):
			body := Body{
				Code:    bizErr.Code,
				Msg:     bizErr.Error(),
				Data:    nil,
				TraceId: GetTraceId(r),
			}
			httpx.OkJsonCtx(r.Context(), w, body)
		case errors.As(err, &unmarshalErr):
			body := Body{
				Code:    http.StatusInternalServerError,
				Msg:     unmarshalErr.Error(),
				Data:    nil,
				TraceId: GetTraceId(r),
			}
			httpx.OkJsonCtx(r.Context(), w, body)
		case errors.As(err, &mysqlErr):
			body := Body{
				Code:    http.StatusInternalServerError,
				Msg:     mysqlErr.Error(),
				Data:    nil,
				TraceId: GetTraceId(r),
			}
			httpx.OkJsonCtx(r.Context(), w, body)
		default:
			body := Body{
				Code:    http.StatusInternalServerError,
				Msg:     err.Error(),
				Data:    nil,
				TraceId: GetTraceId(r),
			}
			httpx.OkJsonCtx(r.Context(), w, body)
		}
		return
	}

	// 2. err为nil的情况，返回成功响应
	body := Body{
		Code:    http.StatusOK,
		Msg:     "successful!",
		Data:    resp,
		TraceId: GetTraceId(r),
	}
	httpx.OkJsonCtx(r.Context(), w, body)
}

// GetTraceId 获取TraceId.
func GetTraceId(r *http.Request) string {
	var traceId string
	spanCtx := trace.SpanContextFromContext(r.Context())
	if spanCtx.HasTraceID() {
		traceId = spanCtx.TraceID().String()
	}

	return traceId
}

// GetLanguage 获取app设置的Language，根据Language返回多语言.
func GetLanguage(r *http.Request) string {
	if len(r.Header["Language"]) > 0 {
		language := r.Header["Language"][0]

		return language
	}

	return "en"
}
