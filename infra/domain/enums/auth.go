package enums

// 验证码场景
const (
	CodeSceneMobileLogin   = "mobile_login"
	CodeSceneEmailLogin    = "email_login"
	CodeSceneRegister      = "register"
	CodeSceneResetPassword = "reset_password"
	CodeSceneBind          = "bind"
)

// IsEmailCodeScene 判断场景是否可用于邮箱验证码。
func IsEmailCodeScene(scene string) bool {
	switch scene {
	case CodeSceneRegister, CodeSceneEmailLogin, CodeSceneResetPassword, CodeSceneBind:
		return true
	default:
		return false
	}
}

// IsMobileCodeScene 判断场景是否可用于手机验证码。
func IsMobileCodeScene(scene string) bool {
	switch scene {
	case CodeSceneMobileLogin, CodeSceneResetPassword, CodeSceneBind:
		return true
	default:
		return false
	}
}

// 认证事件类型 (t_login_log.event_type)
const (
	AuthEventLogin  int64 = 1 // 登录
	AuthEventLogout int64 = 2 // 登出
)

// 认证事件结果 (t_login_log.status)
const (
	AuthResultFail    int64 = 0 // 失败
	AuthResultSuccess int64 = 1 // 成功
)
