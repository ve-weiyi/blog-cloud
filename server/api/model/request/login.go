package request

import (
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/server/utils/valid"
)

type UserReq struct {
	Username string `json:"username" from:"username" example:"admin@qq.com"`
	Password string `json:"password" from:"password" example:"123456"`
	Code     string `json:"code" from:"code" example:""`
}

func (m UserReq) IsValid() error {
	if m.Username == "" || m.Password == "" {
		return fmt.Errorf("用户名或密码不能为空")
	}

	//验证邮箱格式是否正确
	if !valid.IsEmailValid(m.Username) {
		return fmt.Errorf("邮箱格式不正确")
	}

	if len(m.Password) < 6 {
		return fmt.Errorf("密码长度不能小于6位")
	}

	return nil
}

// 用户名只能是邮箱
type UserEmail struct {
	Username string `json:"username" from:"username" example:"admin@qq.com"`
}

func (m UserEmail) IsValid() error {
	//验证邮箱格式是否正确
	if !valid.IsEmailValid(m.Username) {
		return fmt.Errorf("邮箱格式不正确")
	}

	return nil
}

type ResetPasswordReq struct {
	Username string `json:"username" from:"username" example:"admin@qq.com"`
	Password string `json:"password" from:"password" example:"123456"`
	Code     string `json:"code" from:"code" example:""`
}

// Modify password structure
type ChangePasswordReq struct {
	ID          int    `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 旧密码
	NewPassword string `json:"newPassword"` // 新密码
}

type OauthLoginReq struct {
	Platform string `json:"platform" example:""` // 平台
	Code     string `json:"code" example:""`     // 授权码
	State    string `json:"state" example:""`    // 状态
}
