package responsex

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"go.opentelemetry.io/otel/trace"

	"github.com/ve-weiyi/blog-cloud/infra/biz/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizerr"
)

const (
	// HeaderRetryAfter 与限流余量头由产出方（限流中间件、依赖调用方）写入响应，
	// 响应出口只保证不臆造、不覆盖。
	HeaderRetryAfter         = "Retry-After"
	HeaderRateLimitLimit     = "RateLimit-Limit"
	HeaderRateLimitRemaining = "RateLimit-Remaining"
	HeaderRateLimitReset     = "RateLimit-Reset"

	headerWWWAuthenticate = "WWW-Authenticate"
	headerTraceparent     = "traceparent"
	headerCacheControl    = "Cache-Control"

	cacheControlNoStore = "no-store"

	// msgRequestUnparsable 是请求解析失败的 message。它面向开发者供排错，
	// 不透出服务端内部细节；字段级定位由 data.field_violations 承载。
	msgRequestUnparsable = "请求无法解析，请检查请求体格式"
)

// challenge 的四个变体。
//
// challenge 由**标识**决定，不由状态码决定——同为 401，令牌失败、登录凭据失败、
// 与请求签名失败要给出不同的 challenge；RFC 6750 的错误参数只描述 Bearer 令牌，
// 不得套用到其它方案。
var (
	// 未登录 / 登录过期：令牌本身失效，客户端应刷新或登出
	challengeInvalidToken = setHeader(headerWWWAuthenticate, `Bearer realm="blog", error="invalid_token"`)
	// 登录端点凭据失败：请求未出示令牌，不适用 RFC 6750 的错误参数
	challengeBearerOnly = setHeader(headerWWWAuthenticate, `Bearer realm="blog"`)
	// 无权限 / 角色不匹配：凭据有效但权限范围不足（RFC 6750 §3.1）
	challengeInsufficientScope = setHeader(headerWWWAuthenticate, `Bearer realm="blog", error="insufficient_scope"`)
	// 请求签名或时效失败：失败的方案不是 Bearer，用方案自身名
	challengeSignature = setHeader(headerWWWAuthenticate, `Signature realm="blog"`)
)

// rule 是一个标识写入响应时的 HTTP 表达。
type rule struct {
	status  int                       // 出口 HTTP 状态码
	headers func(http.ResponseWriter) // 该标识独有的响应头，可为 nil
}

// rules 是业务错误标识到出口 HTTP 表达的映射。
//
// 映射只在响应出口一处发生——logic 与中间件只产出标识，都不产出状态码。
// 表外标识一律按 500 处理，并由同一包的测试保证表与 bizcode 常量集合一致。
var rules = map[string]rule{
	bizcode.CodeSuccess: {status: http.StatusOK},

	// 429 / 503 的 Retry-After 与 RateLimit-* 由产出方（限流中间件、依赖调用方）在调用本包前
	// 写入 w；出口不臆造：拿不到精确值就省略。
	bizcode.CodeRateLimited: {status: http.StatusTooManyRequests},

	bizcode.CodeRequestSignInvalid: {status: http.StatusUnauthorized, headers: challengeSignature},
	bizcode.CodeUnauthenticated:    {status: http.StatusUnauthorized, headers: challengeInvalidToken},
	bizcode.CodeLoginExpired:       {status: http.StatusUnauthorized, headers: challengeInvalidToken},
	bizcode.CodeCredentialsInvalid: {status: http.StatusUnauthorized, headers: challengeBearerOnly},

	// PARAM_MISSING 是请求头问题，故为 400 而非 422——422 只用于请求内容语义不合法。
	bizcode.CodeInvalidParam:      {status: http.StatusUnprocessableEntity},
	bizcode.CodeParamMissing:      {status: http.StatusBadRequest},
	bizcode.CodeParamFormat:       {status: http.StatusUnprocessableEntity},
	bizcode.CodeParamValueInvalid: {status: http.StatusUnprocessableEntity},
	bizcode.CodeVerifyCodeError:   {status: http.StatusUnprocessableEntity},

	bizcode.CodeAccountDisabled:     {status: http.StatusForbidden},
	bizcode.CodeNoPermission:        {status: http.StatusForbidden, headers: challengeInsufficientScope},
	bizcode.CodeRoleNotMatch:        {status: http.StatusForbidden, headers: challengeInsufficientScope},
	bizcode.CodeOperationNotAllowed: {status: http.StatusForbidden},

	bizcode.CodeResourceNotFound:         {status: http.StatusNotFound},
	bizcode.CodeResourceAlreadyExist:     {status: http.StatusConflict},
	bizcode.CodeResourceStatusNotAllowed: {status: http.StatusConflict},
	bizcode.CodePreconditionFailed:       {status: http.StatusPreconditionFailed},

	bizcode.CodeInternalServerError:  {status: http.StatusInternalServerError},
	bizcode.CodeDatabaseError:        {status: http.StatusInternalServerError},
	bizcode.CodeExternalServiceError: {status: http.StatusBadGateway},
	bizcode.CodeServiceUnavailable:   {status: http.StatusServiceUnavailable},
	bizcode.CodeServiceTimeout:       {status: http.StatusGatewayTimeout},
}

// statusAllowlist 是 ResponseStatus 允许显式声明的状态。
//
// 只收"有状态、无业务标识"的传输层路径：请求解析失败（400）与创建成功（201）。
// 其余状态一律有对应标识，必须走 Response——显式状态出口不得成为绕过标识表的后门。
var statusAllowlist = map[int]bool{
	http.StatusOK:         true,
	http.StatusCreated:    true,
	http.StatusBadRequest: true,
}

// emptyData 是失败响应的 data 取值：字段必须存在，失败时为空对象
var emptyData = struct{}{}

type Body struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
}

// noCodeBody 是无业务标识路径的 body 形状：有 message/data/trace_id，没有 code。
type noCodeBody struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
}

// violationData 是 data 中字段级明细的容器
type violationData struct {
	FieldViolations []bizerr.FieldViolation `json:"field_violations"`
}

// Response 统一封装响应.
// err 类型决定行为:
//   - *bizerr.BizError: body.code 为业务错误标识，HTTP 状态由标识派生
//   - 其他错误:          HTTP 500
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	traceId := GetTraceId(r)
	setTraceparent(w, r)

	if err != nil {
		writeError(r, w, err, traceId)
		return
	}

	if resp == nil {
		resp = emptyData
	}

	writeJson(r, w, http.StatusOK, Body{
		Code:    bizcode.CodeSuccess,
		Message: "",
		Data:    resp,
		TraceId: traceId,
	})
}

// ResponseStatus 写出显式 HTTP 状态的响应，用于"有状态、无业务标识"的传输层路径：
//
//   - 2xx：仍写信封，code 为 SUCCESS——除本入口与 RequestInvalid 外，所有响应都是信封
//   - 非 2xx：写不带 code 的 body，客户端靠响应形状判定，走通用文案分支
//
// data 传 nil 时按该状态的契约给缺省值（非 2xx 给空的 field_violations 数组，该键不可省略）。
// 状态不在白名单内时按内部错误兜底，避免被当作业务状态的后门。
func ResponseStatus(r *http.Request, w http.ResponseWriter, status int, data interface{}) {
	traceId := GetTraceId(r)
	setTraceparent(w, r)

	if !statusAllowlist[status] {
		logx.Errorf("responsex: 状态 %d 不在 ResponseStatus 白名单内，按内部错误兜底；业务错误请用 Response 传标识", status)
		writeJson(r, w, http.StatusInternalServerError, Body{
			Code:    bizcode.CodeInternalServerError,
			Message: "",
			Data:    emptyData,
			TraceId: traceId,
		})
		return
	}

	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		if data == nil {
			data = emptyData
		}

		writeJson(r, w, status, Body{
			Code:    bizcode.CodeSuccess,
			Message: "",
			Data:    data,
			TraceId: traceId,
		})
		return
	}

	if data == nil {
		data = violationData{FieldViolations: []bizerr.FieldViolation{}}
	}

	writeJson(r, w, status, noCodeBody{
		Message: msgRequestUnparsable,
		Data:    data,
		TraceId: traceId,
	})
}

// RequestInvalid 请求读取/解析失败：400、不带 code。
//
// 原始错误进日志与 message 供开发者排错——message 不面向终端用户，客户端在此路径上
// 靠响应形状判定、走通用文案分支。字段级定位由 data.field_violations 承载；go-zero 的
// 解析错误不保证能可靠提取字段名，故给空数组（该键不可省略，无法定位时给空数组）。
func RequestInvalid(r *http.Request, w http.ResponseWriter, err error) {
	traceId := GetTraceId(r)
	setTraceparent(w, r)

	message := msgRequestUnparsable
	if err != nil {
		message = err.Error()
		logx.Errorf("responsex: 请求读取失败: %v", err)
	}

	writeJson(r, w, http.StatusBadRequest, noCodeBody{
		Message: message,
		Data:    violationData{FieldViolations: []bizerr.FieldViolation{}},
		TraceId: traceId,
	})
}

// writeError 把错误翻译成信封与状态码。业务错误标识是 logic 与中间件唯一
// 的表达方式，HTTP 状态码在这里才被派生出来。
func writeError(r *http.Request, w http.ResponseWriter, err error, traceId string) {
	var bizErr *bizerr.BizError
	if errors.As(err, &bizErr) {
		rl, ok := rules[bizErr.Code]
		status := rl.status
		if !ok {
			// 标识未登记映射：按服务端错误兜底，并保留标识本身便于定位
			status = http.StatusInternalServerError
		}

		if rl.headers != nil {
			rl.headers(w)
		}

		writeJson(r, w, status, Body{
			Code:    bizErr.Code,
			Message: bizErr.Error(),
			Data:    errorData(bizErr),
			TraceId: traceId,
		})
		return
	}

	writeJson(r, w, http.StatusInternalServerError, Body{
		Code:    bizcode.CodeInternalServerError,
		Message: err.Error(),
		Data:    emptyData,
		TraceId: traceId,
	})
}

// errorData 给出失败响应的 data：有字段级明细时带上，否则为空对象
func errorData(bizErr *bizerr.BizError) interface{} {
	if len(bizErr.Violations) == 0 {
		return emptyData
	}

	return violationData{FieldViolations: bizErr.Violations}
}

// writeJson 写出响应，并补与该响应性质绑定的头：
// 所有非 2xx 都不得被中间层或浏览器缓存。
func writeJson(r *http.Request, w http.ResponseWriter, status int, body interface{}) {
	if status < http.StatusOK || status >= http.StatusMultipleChoices {
		w.Header().Set(headerCacheControl, cacheControlNoStore)
	}

	httpx.WriteJsonCtx(r.Context(), w, status, body)
}

// setHeader 返回一个只在响应上写一个头的函数，供 rules 表按标识挂头。
func setHeader(name string, value string) func(http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		w.Header().Set(name, value)
	}
}

// setTraceparent 让网关、反代与 APM 不必解析 body 即可关联链路。
// 与信封的 trace_id 同源，头名用 W3C Trace Context 标准以对接 OpenTelemetry。
func setTraceparent(w http.ResponseWriter, r *http.Request) {
	spanCtx := trace.SpanContextFromContext(r.Context())
	if !spanCtx.IsValid() {
		return
	}

	w.Header().Set(headerTraceparent, fmt.Sprintf("00-%s-%s-%02x",
		spanCtx.TraceID(), spanCtx.SpanID(), spanCtx.TraceFlags()))
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
