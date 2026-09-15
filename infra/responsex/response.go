package responsex

import (
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go.opentelemetry.io/otel/trace"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/biz/httperr"
)

type Body struct {
	Code        int64       `json:"code"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	EncryptData string      `json:"encrypt_data,omitempty"`
	TraceId     string      `json:"trace_id"`
}

// Response 统一封装响应.
// err 类型决定行为:
//   - *httperr.HttpError: Code 直接作为 HTTP 响应状态码
//   - *bizerr.BizError:   HTTP 200，body.code 为业务状态码（保持向后兼容）
//   - 其他错误:            HTTP 500
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		var httpErr *httperr.HttpError
		var bizErr *bizerr.BizError

		switch {
		case errors.As(err, &httpErr):
			httpx.WriteJsonCtx(r.Context(), w, int(httpErr.Code), Body{
				Code:    httpErr.Code,
				Message: httpErr.Message,
				TraceId: GetTraceId(r),
			})
		case errors.As(err, &bizErr):
			httpx.OkJsonCtx(r.Context(), w, Body{
				Code:    bizErr.Code,
				Message: bizErr.Error(),
				TraceId: GetTraceId(r),
			})
		default:
			httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, Body{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				TraceId: GetTraceId(r),
			})
		}
		return
	}

	httpx.OkJsonCtx(r.Context(), w, Body{
		Code:    http.StatusOK,
		Message: "successful!",
		Data:    resp,
		TraceId: GetTraceId(r),
	})
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
