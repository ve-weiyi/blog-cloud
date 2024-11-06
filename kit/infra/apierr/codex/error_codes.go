package codex

// 错误码定义规则
const (
	// 3xx 请求频繁
	CodeTooManyRequests = 300 // 操作频繁

	// 4xx 请求参数错误
	CodeInvalidParam      = 400 // 请求参数错误
	CodeUserUnLogin       = 401 // 用户未登录
	CodeUserDisabled      = 402 // 用户账号错误
	CodeUserNotPermission = 403 // 用户无权限
	CodeUserNotExist      = 404 // 用户不存在
	CodeUserAlreadyExist  = 405 // 用户已存在
	CodeUserPasswordError = 406 // 用户密码错误

	// 5xx 服务器错误
	CodeInternalServerError = 500 // 内部错误
	CodeSqlQueryError       = 501 //数据库查询错误
	CodeCaptchaVerify       = 502 //验证码错误
)
