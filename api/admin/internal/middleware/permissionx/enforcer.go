package permissionx

type Enforcer interface {
	// 加载权限
	LoadPolicy() error
	// 重新加载权限
	ReloadPolicy() error
	// 验证用户是否有权限
	Enforce(user string, resource string, action string) (bool, error)
}
