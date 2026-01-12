package httperr

// 定义HTTP状态码常量
const (
	// 1xx Informational（信息性状态码）
	CodeContinue           = 100 // 继续，服务器收到请求，客户端应继续发送请求的其余部分
	CodeSwitchingProtocols = 101 // 切换协议，服务器已经理解客户端的请求，并将通过Upgrade消息头通知客户端切换协议
	CodeProcessing         = 102 // 处理中，服务器正在处理请求，但尚未完成

	// 2xx Success（成功状态码）
	CodeOK                   = 200 // 请求成功，服务器成功处理了请求
	CodeCreated              = 201 // 已创建，请求已经被实现，而且有一个新的资源已经依据请求的需要而建立
	CodeAccepted             = 202 // 已接受，服务器已接受请求，但尚未处理
	CodeNonAuthoritativeInfo = 203 // 非权威信息，服务器已成功处理了请求，但返回的信息可能来自另一来源
	CodeNoContent            = 204 // 无内容，服务器成功处理了请求，但没有返回任何内容
	CodeResetContent         = 205 // 重置内容，服务器成功处理了请求，且重置了当前页面的内容
	CodePartialContent       = 206 // 部分内容，服务器成功处理了部分GET请求
	CodeMultiStatus          = 207 // 多状态，服务器已经执行了一部分GET请求
	CodeAlreadyReported      = 208 // 已报告，服务器已经在上一个请求的响应中发出了对该状态码的说明
	CodeIMUsed               = 226 // IM使用，服务器已经满足了对资源的请求

	// 3xx Redirection（重定向状态码）
	CodeMultipleChoices   = 300 // 多种选择，请求的资源可包括多个位置
	CodeMovedPermanently  = 301 // 永久移动，请求的资源已被永久移动到新位置
	CodeFound             = 302 // 已找到，请求的资源临时从不同的URI响应请求
	CodeSeeOther          = 303 // 查看其他位置，服务器返回此响应时，会提供一个URI以参考
	CodeNotModified       = 304 // 未修改，自从上次请求后，请求的资源未修改
	CodeUseProxy          = 305 // 使用代理，请求者只能使用代理访问请求的资源
	CodeTemporaryRedirect = 307 // 临时重定向，请求的资源临时从不同的URI响应请求
	CodePermanentRedirect = 308 // 永久重定向，请求的资源已被永久移动到新位置

	// 4xx Client Errors（客户端错误状态码）
	CodeBadRequest                  = 400 // 请求错误，服务器不理解请求的语法
	CodeUnauthorized                = 401 // 未授权，请求要求身份验证
	CodePaymentRequired             = 402 // 需要付款，保留供将来使用
	CodeForbidden                   = 403 // 禁止访问，服务器已经理解请求，但是拒绝执行它
	CodeNotFound                    = 404 // 未找到，服务器找不到请求的资源
	CodeMethodNotAllowed            = 405 // 方法不允许，请求中的方法被禁止
	CodeNotAcceptable               = 406 // 不可接受，服务器无法根据请求的内容特性完成请求
	CodeProxyAuthRequired           = 407 // 代理认证要求，客户端必须先使用代理认证自身
	CodeRequestTimeout              = 408 // 请求超时，服务器等待请求时发生超时
	CodeConflict                    = 409 // 冲突，由于请求中的冲突，请求无法被完成
	CodeGone                        = 410 // 资源不可用，请求的资源在服务器上已经不存在
	CodeLengthRequired              = 411 // 需要内容长度，服务器拒绝在没有定义 Content-Length 头的情况下接受请求
	CodePreconditionFailed          = 412 // 先决条件失败，请求头中指定的一些前提条件失败
	CodePayloadTooLarge             = 413 // 请求实体过大，服务器无法处理此请求
	CodeURITooLong                  = 414 // 请求的URI过长，服务器无法处理
	CodeUnsupportedMediaType        = 415 // 不支持的媒体类型，服务器无法处理请求使用的媒体类型
	CodeRangeNotSatisfiable         = 416 // 范围不符合要求，客户端请求的范围无效
	CodeExpectationFailed           = 417 // 期望失败，服务器未满足"期望"请求标头字段的要求
	CodeMisdirectedRequest          = 421 // 误导的请求，服务器不会按照无效或不可理解的消息头字段的定义来处理请求
	CodeUnprocessableEntity         = 422 // 无法处理的实体，请求格式正确，但由于含有语义错误，无法响应
	CodeLocked                      = 423 // 已锁定，当前资源被锁定
	CodeFailedDependency            = 424 // 依赖失败，由于之前的请求失败，导致请求失败
	CodeTooEarly                    = 425 // 过早，服务器不愿意冒风险处理请求，可能造成重放攻击
	CodeUpgradeRequired             = 426 // 需要升级，客户端应切换到TLS/1.0
	CodePreconditionRequired        = 428 // 先决条件要求，原始服务器要求该请求是有条件的
	CodeTooManyRequests             = 429 // 请求过多，用户在给定的时间内发送了太多请求
	CodeRequestHeaderFieldsTooLarge = 431 // 请求头字段太大，服务器不愿意处理请求，因为它的请求头字段太大
	CodeUnavailableForLegalReasons  = 451 // 由于法律原因不可用，用户请求非法资源，例如：由政府审查的网页

	// 5xx Server Errors（服务器错误状态码）
	CodeInternalServerError           = 500 // 服务器内部错误，服务器遇到错误，无法完成请求
	CodeNotImplemented                = 501 // 未实现，服务器不支持当前请求所需要的某个功能
	CodeBadGateway                    = 502 // 错误的网关，服务器作为网关或代理，从上游服务器收到无效响应
	CodeServiceUnavailable            = 503 // 服务不可用，服务器目前无法使用（由于超载或停机维护）
	CodeGatewayTimeout                = 504 // 网关超时，服务器作为网关或代理，但是没有及时从上游服务器收到请求
	CodeHTTPVersionNotSupported       = 505 // HTTP版本不受支持，服务器不支持请求中所用的HTTP协议版本
	CodeVariantAlsoNegotiates         = 506 // 内容协商失败，服务器无法完成对资源的请求
	CodeInsufficientStorage           = 507 // 存储空间不足，服务器无法完成存储请求所需的内容
	CodeLoopDetected                  = 508 // 检测到无限循环，服务器检测到无限循环
	CodeNotExtended                   = 510 // 未扩展，获取资源所需的策略没有被满足
	CodeNetworkAuthenticationRequired = 511 // 网络认证要求，客户端需要进行身份验证才能获得网络访问权限
)
