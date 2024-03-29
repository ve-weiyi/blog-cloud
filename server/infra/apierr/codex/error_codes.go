package codex

// 错误码定义规则
const (
	// 3xx 用户行为错误
	CodeUserUnLogin       = 401 // 用户未登录
	CodeUserDisabled      = 402 // 用户已被禁用
	CodeUserNotPermission = 403 // 用户无权限
	CodeUserNotExist      = 404 // 用户不存在
	CodeUserAlreadyExist  = 405 // 用户已存在
	CodeUserPasswordError = 406 // 用户密码错误

	// 4xx 请求参数错误
	CodeInvalidParam    = 400 // 请求参数错误
	CodeTooManyRequests = 401 // 操作频繁

	// 5xx 服务器错误
	CodeInternalServerError = 500 // 内部错误
	CodeCaptchaVerifyError  = 501 // 验证码错误
	CodeSqlQueryError       = 502 //数据库查询错误
)
