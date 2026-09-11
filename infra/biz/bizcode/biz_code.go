package bizcode

const (
	CodeSuccess = 200
)

// 3xx 流控/安全
const (
	CodeTooManyRequests    = 300 // 请求过于频繁
	CodeRateLimited        = 301 // 请求被限流
	CodeCaptchaError       = 302 // 验证码校验失败
	CodeRequestSignInvalid = 303 // 请求签名无效
)

// 4xx 客户端错误

// 400-409 请求参数
const (
	CodeInvalidParam      = 400 // 参数错误
	CodeParamMissing      = 401 // 参数缺失
	CodeParamFormat       = 402 // 参数格式错误
	CodeParamValueInvalid = 403 // 参数值不允许
)

// 410-419 身份认证
const (
	CodeUnauthenticated = 410 // 未登录
	CodeLoginExpired    = 411 // 登录已过期
	CodePasswordError   = 412 // 密码错误
	CodeAccountDisabled = 413 // 账号已被禁用
	CodeVerifyCodeError = 414 // 验证码错误
)

// 420-429 权限授权
const (
	CodeNoPermission = 420 // 无操作权限
	CodeRoleNotMatch = 421 // 角色不匹配
)

// 430-449 业务规则
const (
	CodeResourceNotFound         = 430 // 资源不存在
	CodeResourceAlreadyExist     = 431 // 资源已存在
	CodeResourceStatusNotAllowed = 432 // 资源状态不允许当前操作
	CodeOperationNotAllowed      = 433 // 操作不被允许
)

// 5xx 服务端错误
const (
	CodeInternalServerError  = 500 // 服务器内部错误
	CodeDatabaseError        = 501 // 数据库操作失败
	CodeExternalServiceError = 502 // 外部服务调用失败
	CodeServiceTimeout       = 503 // 服务超时
)
