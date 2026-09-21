// Package bizcode 定义本项目的业务错误标识。
//
// 标识自描述、只增不改、不复用；同一语义只有一个标识，可跨模块复用。
//
// 标识与 HTTP 状态码的映射在响应出口（infra/responsex）逐标识定义，
// 本包不承载 HTTP 语义。
package bizcode

const (
	CodeSuccess = "SUCCESS" // 成功
)

// 流控与安全
const (
	CodeRateLimited        = "RATE_LIMIT_EXCEEDED"       // 请求被限流
	CodeRequestSignInvalid = "REQUEST_SIGNATURE_INVALID" // 请求签名无效
)

// 请求参数
const (
	CodeInvalidParam      = "PARAM_INVALID"           // 参数错误（通用兜底）
	CodeParamMissing      = "PARAM_MISSING"           // 参数缺失
	CodeParamFormat       = "PARAM_FORMAT_INVALID"    // 参数格式错误
	CodeParamValueInvalid = "PARAM_VALUE_NOT_ALLOWED" // 参数值不允许
)

// 身份认证
const (
	CodeUnauthenticated    = "UNAUTHENTICATED"     // 未登录
	CodeLoginExpired       = "LOGIN_EXPIRED"       // 登录已过期
	CodeCredentialsInvalid = "CREDENTIALS_INVALID" // 凭据无效（账号不存在与密码错误合并返回）
	CodeAccountDisabled    = "ACCOUNT_DISABLED"    // 账号已被禁用
	CodeVerifyCodeError    = "CAPTCHA_INCORRECT"   // 验证码错误
)

// 权限授权
const (
	CodeNoPermission = "PERMISSION_DENIED" // 无操作权限
	CodeRoleNotMatch = "ROLE_MISMATCH"     // 角色不匹配
)

// 业务规则
const (
	CodeResourceNotFound         = "RESOURCE_NOT_FOUND"         // 资源不存在
	CodeResourceAlreadyExist     = "RESOURCE_ALREADY_EXISTS"    // 资源已存在
	CodeResourceStatusNotAllowed = "RESOURCE_STATE_NOT_ALLOWED" // 资源状态不允许当前操作
	CodePreconditionFailed       = "PRECONDITION_FAILED"        // 条件请求的前提不成立（并发编辑冲突）
	CodeOperationNotAllowed      = "OPERATION_NOT_ALLOWED"      // 操作不被允许
)

// 服务端错误
const (
	CodeInternalServerError  = "INTERNAL_ERROR"         // 服务器内部错误
	CodeDatabaseError        = "DATABASE_ERROR"         // 数据库操作失败
	CodeExternalServiceError = "EXTERNAL_SERVICE_ERROR" // 外部服务调用失败
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"    // 依赖暂时不可用
	CodeServiceTimeout       = "SERVICE_TIMEOUT"        // 服务超时
)

// All 列出全部业务错误标识，作为集合的完整枚举。
//
// **新增标识时必须同时更新本切片**，否则 responsex 的映射表一致性测试
// 无法发现遗漏（新增标识必须同时确定产出层与 HTTP 映射）。
var All = []string{
	CodeSuccess,

	CodeRateLimited,
	CodeRequestSignInvalid,

	CodeInvalidParam,
	CodeParamMissing,
	CodeParamFormat,
	CodeParamValueInvalid,

	CodeUnauthenticated,
	CodeLoginExpired,
	CodeCredentialsInvalid,
	CodeAccountDisabled,
	CodeVerifyCodeError,

	CodeNoPermission,
	CodeRoleNotMatch,

	CodeResourceNotFound,
	CodeResourceAlreadyExist,
	CodeResourceStatusNotAllowed,
	CodePreconditionFailed,
	CodeOperationNotAllowed,

	CodeInternalServerError,
	CodeDatabaseError,
	CodeExternalServiceError,
	CodeServiceUnavailable,
	CodeServiceTimeout,
}
