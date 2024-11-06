package rbacx

type IApiPolicy interface {
	HasPermission(roles []string) bool
	Module() string
	Desc() string
	Disable() bool
	Traceable() bool
}

type ApiPolicy struct {
	module string
	desc   string

	disable   bool
	traceable bool
	roles     []string
}

// 判断是否有权限访问
func (m *ApiPolicy) HasPermission(roles []string) bool {
	if len(m.roles) == 0 {
		return true
	}

	for _, v := range roles {
		for _, r := range m.roles {
			if v == r {
				return true
			}
		}
	}

	return false
}

func (m *ApiPolicy) Module() string {
	return m.module
}

func (m *ApiPolicy) Desc() string {
	return m.desc
}

func (m *ApiPolicy) Disable() bool {
	return m.disable
}

func (m *ApiPolicy) Traceable() bool {
	return m.traceable
}

var _ IApiPolicy = (*ApiPolicy)(nil)
