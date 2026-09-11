package httperr

// 定义ERROR常量
// 1xx Informational（信息性状态码）
var (
	ErrorContinue           = NewHttpError(CodeContinue, "继续")
	ErrorSwitchingProtocols = NewHttpError(CodeSwitchingProtocols, "切换协议")
	ErrorProcessing         = NewHttpError(CodeProcessing, "处理中")
)

// 2xx Success（成功状态码）
var (
	ErrorOK                   = NewHttpError(CodeOK, "请求成功")
	ErrorCreated              = NewHttpError(CodeCreated, "已创建")
	ErrorAccepted             = NewHttpError(CodeAccepted, "已接受")
	ErrorNonAuthoritativeInfo = NewHttpError(CodeNonAuthoritativeInfo, "非权威信息")
	ErrorNoContent            = NewHttpError(CodeNoContent, "无内容")
	ErrorResetContent         = NewHttpError(CodeResetContent, "重置内容")
	ErrorPartialContent       = NewHttpError(CodePartialContent, "部分内容")
	ErrorMultiStatus          = NewHttpError(CodeMultiStatus, "多状态")
	ErrorAlreadyReported      = NewHttpError(CodeAlreadyReported, "已报告")
	ErrorIMUsed               = NewHttpError(CodeIMUsed, "使用了IM")
)

// 3xx Redirection（重定向状态码）
var (
	ErrorMultipleChoices   = NewHttpError(CodeMultipleChoices, "多种选择")
	ErrorMovedPermanently  = NewHttpError(CodeMovedPermanently, "永久移动")
	ErrorFound             = NewHttpError(CodeFound, "已找到")
	ErrorSeeOther          = NewHttpError(CodeSeeOther, "查看其他位置")
	ErrorNotModified       = NewHttpError(CodeNotModified, "未修改")
	ErrorUseProxy          = NewHttpError(CodeUseProxy, "使用代理")
	ErrorTemporaryRedirect = NewHttpError(CodeTemporaryRedirect, "临时重定向")
	ErrorPermanentRedirect = NewHttpError(CodePermanentRedirect, "永久重定向")
)

// 4xx Client Errors（客户端错误状态码）
var (
	ErrorBadRequest                  = NewHttpError(CodeBadRequest, "请求错误")
	ErrorUnauthorized                = NewHttpError(CodeUnauthorized, "未授权")
	ErrorPaymentRequired             = NewHttpError(CodePaymentRequired, "需要付款")
	ErrorForbidden                   = NewHttpError(CodeForbidden, "禁止访问")
	ErrorNotFound                    = NewHttpError(CodeNotFound, "未找到")
	ErrorMethodNotAllowed            = NewHttpError(CodeMethodNotAllowed, "方法不允许")
	ErrorNotAcceptable               = NewHttpError(CodeNotAcceptable, "不可接受")
	ErrorProxyAuthRequired           = NewHttpError(CodeProxyAuthRequired, "代理认证要求")
	ErrorRequestTimeout              = NewHttpError(CodeRequestTimeout, "请求超时")
	ErrorConflict                    = NewHttpError(CodeConflict, "冲突")
	ErrorGone                        = NewHttpError(CodeGone, "资源不可用")
	ErrorLengthRequired              = NewHttpError(CodeLengthRequired, "需要内容长度")
	ErrorPreconditionFailed          = NewHttpError(CodePreconditionFailed, "先决条件失败")
	ErrorPayloadTooLarge             = NewHttpError(CodePayloadTooLarge, "请求实体过大")
	ErrorURITooLong                  = NewHttpError(CodeURITooLong, "请求的URI过长")
	ErrorUnsupportedMediaType        = NewHttpError(CodeUnsupportedMediaType, "不支持的媒体类型")
	ErrorRangeNotSatisfiable         = NewHttpError(CodeRangeNotSatisfiable, "范围不符合要求")
	ErrorExpectationFailed           = NewHttpError(CodeExpectationFailed, "期望失败")
	ErrorMisdirectedRequest          = NewHttpError(CodeMisdirectedRequest, "误导的请求")
	ErrorUnprocessableEntity         = NewHttpError(CodeUnprocessableEntity, "无法处理的实体")
	ErrorLocked                      = NewHttpError(CodeLocked, "已锁定")
	ErrorFailedDependency            = NewHttpError(CodeFailedDependency, "依赖失败")
	ErrorTooEarly                    = NewHttpError(CodeTooEarly, "过早")
	ErrorUpgradeRequired             = NewHttpError(CodeUpgradeRequired, "需要升级")
	ErrorPreconditionRequired        = NewHttpError(CodePreconditionRequired, "先决条件要求")
	ErrorTooManyRequests             = NewHttpError(CodeTooManyRequests, "请求过多")
	ErrorRequestHeaderFieldsTooLarge = NewHttpError(CodeRequestHeaderFieldsTooLarge, "请求头字段太大")
	ErrorUnavailableForLegalReasons  = NewHttpError(CodeUnavailableForLegalReasons, "由于法律原因不可用")
)

// 5xx Server Errors（服务器错误状态码）
var (
	ErrorInternalServerError           = NewHttpError(CodeInternalServerError, "服务器内部错误")
	ErrorNotImplemented                = NewHttpError(CodeNotImplemented, "未实现")
	ErrorBadGateway                    = NewHttpError(CodeBadGateway, "错误的网关")
	ErrorServiceUnavailable            = NewHttpError(CodeServiceUnavailable, "服务不可用")
	ErrorGatewayTimeout                = NewHttpError(CodeGatewayTimeout, "网关超时")
	ErrorHTTPVersionNotSupported       = NewHttpError(CodeHTTPVersionNotSupported, "不支持的HTTP版本")
	ErrorVariantAlsoNegotiates         = NewHttpError(CodeVariantAlsoNegotiates, "变体也协商")
	ErrorInsufficientStorage           = NewHttpError(CodeInsufficientStorage, "存储不足")
	ErrorLoopDetected                  = NewHttpError(CodeLoopDetected, "检测到循环")
	ErrorNotExtended                   = NewHttpError(CodeNotExtended, "未扩展")
	ErrorNetworkAuthenticationRequired = NewHttpError(CodeNetworkAuthenticationRequired, "需要网络认证")
)
