package bizerr

var (
	ErrorBadRequest = NewBizError(400, "请求错误")
	//ErrorInvalidParam    = NewBizError(4001, "请求参数错误")
	//ErrorFrequentRequest = NewBizError(4002, "操作频繁,请在10分钟后再试")
	//
	//ErrorUnauthorized = NewBizError(401, "未授权")
	//
	//ErrorUserDisabled      = NewBizError(4012, "用户已被禁用")
	//ErrorUserNotPermission = NewBizError(4013, "用户无操作权限")
	//ErrorUserNotExist      = NewBizError(4014, "用户不存在")
	//ErrorUserAlreadyExist  = NewBizError(4015, "用户已存在")
	//ErrorUserPasswordError = NewBizError(4016, "用户密码错误")
	//
	//ErrorInternalServerError = NewBizError(500, "服务器内部错误")
	//ErrorCaptchaVerify       = NewBizError(5001, "验证码错误")
	//ErrorSqlQueryError       = NewBizError(5002, "数据库查询错误")
)
