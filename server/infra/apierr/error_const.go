package apierr

import (
	"github.com/ve-weiyi/ve-blog-golang/server/infra/apierr/codex"
)

var (
	ErrorInvalidParam    = NewApiError(codex.CodeInvalidParam, "请求参数错误")
	ErrorCaptchaVerify   = NewApiError(codex.CodeCaptchaVerifyError, "验证码错误")
	ErrorFrequentRequest = NewApiError(codex.CodeTooManyRequests, "操作频繁,请在10分钟后再试")

	ErrorUserUnLogin       = NewApiError(codex.CodeUserUnLogin, "用户未登录")
	ErrorUserDisabled      = NewApiError(codex.CodeUserDisabled, "用户已被禁用")
	ErrorUserNotPermission = NewApiError(codex.CodeUserNotPermission, "用户无操作权限")
	ErrorUserNotExist      = NewApiError(codex.CodeUserNotExist, "用户不存在")
	ErrorUserAlreadyExist  = NewApiError(codex.CodeUserAlreadyExist, "用户已存在")
	ErrorUserPasswordError = NewApiError(codex.CodeUserPasswordError, "用户密码错误")

	ErrorInternalServerError = NewApiError(codex.CodeInternalServerError, "服务器内部错误")
	ErrorSqlQueryError       = NewApiError(codex.CodeSqlQueryError, "数据库查询错误")
)
