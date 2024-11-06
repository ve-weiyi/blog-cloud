package apierr

var (
	ErrorBadRequest      = NewApiError(400, "请求错误")
	ErrorInvalidParam    = NewApiError(4001, "请求参数错误")
	ErrorFrequentRequest = NewApiError(4002, "操作频繁,请在10分钟后再试")

	ErrorUnauthorized = NewApiError(401, "未授权")

	ErrorUserDisabled      = NewApiError(4012, "用户已被禁用")
	ErrorUserNotPermission = NewApiError(4013, "用户无操作权限")
	ErrorUserNotExist      = NewApiError(4014, "用户不存在")
	ErrorUserAlreadyExist  = NewApiError(4015, "用户已存在")
	ErrorUserPasswordError = NewApiError(4016, "用户密码错误")

	ErrorInternalServerError = NewApiError(500, "服务器内部错误")
	ErrorCaptchaVerify       = NewApiError(5001, "验证码错误")
	ErrorSqlQueryError       = NewApiError(5002, "数据库查询错误")
)
