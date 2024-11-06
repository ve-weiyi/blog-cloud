package httperr

import (
	"github.com/ve-weiyi/ve-blog-golang/kit/infra/apierr"
)

// 定义ERROR常量
// 1xx Informational（信息性状态码）
var (
	ErrorContinue           = apierr.NewApiError(CodeContinue, "继续")
	ErrorSwitchingProtocols = apierr.NewApiError(CodeSwitchingProtocols, "切换协议")
	ErrorProcessing         = apierr.NewApiError(CodeProcessing, "处理中")
)

// 2xx Success（成功状态码）
var (
	ErrorOK                   = apierr.NewApiError(CodeOK, "请求成功")
	ErrorCreated              = apierr.NewApiError(CodeCreated, "已创建")
	ErrorAccepted             = apierr.NewApiError(CodeAccepted, "已接受")
	ErrorNonAuthoritativeInfo = apierr.NewApiError(CodeNonAuthoritativeInfo, "非权威信息")
	ErrorNoContent            = apierr.NewApiError(CodeNoContent, "无内容")
	ErrorResetContent         = apierr.NewApiError(CodeResetContent, "重置内容")
	ErrorPartialContent       = apierr.NewApiError(CodePartialContent, "部分内容")
	ErrorMultiStatus          = apierr.NewApiError(CodeMultiStatus, "多状态")
	ErrorAlreadyReported      = apierr.NewApiError(CodeAlreadyReported, "已报告")
	ErrorIMUsed               = apierr.NewApiError(CodeIMUsed, "使用了IM")
)

// 3xx Redirection（重定向状态码）
var (
	ErrorMultipleChoices   = apierr.NewApiError(CodeMultipleChoices, "多种选择")
	ErrorMovedPermanently  = apierr.NewApiError(CodeMovedPermanently, "永久移动")
	ErrorFound             = apierr.NewApiError(CodeFound, "已找到")
	ErrorSeeOther          = apierr.NewApiError(CodeSeeOther, "查看其他位置")
	ErrorNotModified       = apierr.NewApiError(CodeNotModified, "未修改")
	ErrorUseProxy          = apierr.NewApiError(CodeUseProxy, "使用代理")
	ErrorTemporaryRedirect = apierr.NewApiError(CodeTemporaryRedirect, "临时重定向")
	ErrorPermanentRedirect = apierr.NewApiError(CodePermanentRedirect, "永久重定向")
)

// 4xx Client Errors（客户端错误状态码）
var (
	ErrorBadRequest                  = apierr.NewApiError(CodeBadRequest, "请求错误")
	ErrorUnauthorized                = apierr.NewApiError(CodeUnauthorized, "未授权")
	ErrorPaymentRequired             = apierr.NewApiError(CodePaymentRequired, "需要付款")
	ErrorForbidden                   = apierr.NewApiError(CodeForbidden, "禁止访问")
	ErrorNotFound                    = apierr.NewApiError(CodeNotFound, "未找到")
	ErrorMethodNotAllowed            = apierr.NewApiError(CodeMethodNotAllowed, "方法不允许")
	ErrorNotAcceptable               = apierr.NewApiError(CodeNotAcceptable, "不可接受")
	ErrorProxyAuthRequired           = apierr.NewApiError(CodeProxyAuthRequired, "代理认证要求")
	ErrorRequestTimeout              = apierr.NewApiError(CodeRequestTimeout, "请求超时")
	ErrorConflict                    = apierr.NewApiError(CodeConflict, "冲突")
	ErrorGone                        = apierr.NewApiError(CodeGone, "资源不可用")
	ErrorLengthRequired              = apierr.NewApiError(CodeLengthRequired, "需要内容长度")
	ErrorPreconditionFailed          = apierr.NewApiError(CodePreconditionFailed, "先决条件失败")
	ErrorPayloadTooLarge             = apierr.NewApiError(CodePayloadTooLarge, "请求实体过大")
	ErrorURITooLong                  = apierr.NewApiError(CodeURITooLong, "请求的URI过长")
	ErrorUnsupportedMediaType        = apierr.NewApiError(CodeUnsupportedMediaType, "不支持的媒体类型")
	ErrorRangeNotSatisfiable         = apierr.NewApiError(CodeRangeNotSatisfiable, "范围不符合要求")
	ErrorExpectationFailed           = apierr.NewApiError(CodeExpectationFailed, "期望失败")
	ErrorMisdirectedRequest          = apierr.NewApiError(CodeMisdirectedRequest, "误导的请求")
	ErrorUnprocessableEntity         = apierr.NewApiError(CodeUnprocessableEntity, "无法处理的实体")
	ErrorLocked                      = apierr.NewApiError(CodeLocked, "已锁定")
	ErrorFailedDependency            = apierr.NewApiError(CodeFailedDependency, "依赖失败")
	ErrorTooEarly                    = apierr.NewApiError(CodeTooEarly, "过早")
	ErrorUpgradeRequired             = apierr.NewApiError(CodeUpgradeRequired, "需要升级")
	ErrorPreconditionRequired        = apierr.NewApiError(CodePreconditionRequired, "先决条件要求")
	ErrorTooManyRequests             = apierr.NewApiError(CodeTooManyRequests, "请求过多")
	ErrorRequestHeaderFieldsTooLarge = apierr.NewApiError(CodeRequestHeaderFieldsTooLarge, "请求头字段太大")
	ErrorUnavailableForLegalReasons  = apierr.NewApiError(CodeUnavailableForLegalReasons, "由于法律原因不可用")
)

// 5xx Server Errors（服务器错误状态码）
var (
	ErrorInternalServerError           = apierr.NewApiError(CodeInternalServerError, "服务器内部错误")
	ErrorNotImplemented                = apierr.NewApiError(CodeNotImplemented, "未实现")
	ErrorBadGateway                    = apierr.NewApiError(CodeBadGateway, "错误的网关")
	ErrorServiceUnavailable            = apierr.NewApiError(CodeServiceUnavailable, "服务不可用")
	ErrorGatewayTimeout                = apierr.NewApiError(CodeGatewayTimeout, "网关超时")
	ErrorHTTPVersionNotSupported       = apierr.NewApiError(CodeHTTPVersionNotSupported, "不支持的HTTP版本")
	ErrorVariantAlsoNegotiates         = apierr.NewApiError(CodeVariantAlsoNegotiates, "变体也协商")
	ErrorInsufficientStorage           = apierr.NewApiError(CodeInsufficientStorage, "存储不足")
	ErrorLoopDetected                  = apierr.NewApiError(CodeLoopDetected, "检测到循环")
	ErrorNotExtended                   = apierr.NewApiError(CodeNotExtended, "未扩展")
	ErrorNetworkAuthenticationRequired = apierr.NewApiError(CodeNetworkAuthenticationRequired, "需要网络认证")
)
