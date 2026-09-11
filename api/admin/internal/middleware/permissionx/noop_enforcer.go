package permissionx

import "github.com/zeromicro/go-zero/core/logx"

var _ Enforcer = &NoopEnforcer{}

// NoopEnforcer 允许所有操作
type NoopEnforcer struct{}

func NewNoopEnforcer() *NoopEnforcer {
	return &NoopEnforcer{}
}

func (m *NoopEnforcer) ReloadPolicy() error {
	return m.LoadPolicy()
}

func (m *NoopEnforcer) LoadPolicy() error {
	logx.Info("Reloading permissions...")
	return nil
}

func (m *NoopEnforcer) Enforce(user string, resource string, action string) (bool, error) {
	return false, nil
}
