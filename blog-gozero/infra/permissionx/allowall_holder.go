package permissionx

import (
	"github.com/zeromicro/go-zero/core/logx"
)

// 允许所有操作
type AllowAllHolder struct {
}

func NewAllowAllHolder() *AllowAllHolder {
	return &AllowAllHolder{}
}

func (m *AllowAllHolder) ReloadPolicy() error {
	return m.LoadPolicy()
}

func (m *AllowAllHolder) LoadPolicy() error {
	logx.Info("Reloading permissions...")
	return nil
}

func (m *AllowAllHolder) Enforce(user string, resource string, action string) error {
	// todo
	return nil
}

func (m *AllowAllHolder) dynamicLoadUserRoles(userId string) error {
	// todo
	return nil
}

func (m *AllowAllHolder) enforce(user string, resource string, action string) (bool, error) {
	// todo
	return true, nil
}
