package permissionx

import "github.com/zeromicro/go-zero/core/logx"

var _ Enforcer = &NoopEnforcer{}

// NoopEnforcer 空实现：不加载策略、放行所有操作。
// 仅用于本地调试或关闭鉴权的场景，生产环境禁止接入 —— 接入即等于不做任何鉴权。
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
	return true, nil
}
