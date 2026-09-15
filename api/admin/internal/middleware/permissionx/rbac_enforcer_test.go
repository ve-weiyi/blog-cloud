package permissionx

import (
	"testing"

	"github.com/mikespook/gorbac/v2"
)

// registered 是「已注册接口」全集，用于填充 allPerms。
// Enforce 对未注册的接口按设计放行（见 rbac_enforcer.go 的 unregistered 分支），
// 因此被拒绝的用例必须先把接口登记进来，否则会走放行分支而测不到鉴权。
func newTestEnforcer(userRoles map[string][]string, rolePerms map[string][]string, registered []string) *RbacEnforcer {
	rbac := gorbac.New()
	for role, perms := range rolePerms {
		r := gorbac.NewStdRole(role)
		for _, perm := range perms {
			_ = r.Assign(apiPermission(perm))
		}
		_ = rbac.Add(r)
	}

	allPerms := make(map[string]struct{}, len(registered))
	for _, perm := range registered {
		allPerms[perm] = struct{}{}
	}

	return &RbacEnforcer{
		rbac:         rbac,
		userRoles:    userRoles,
		allPerms:     allPerms,
		policyLoaded: true,
	}
}

func TestEnforce(t *testing.T) {
	e := newTestEnforcer(
		map[string][]string{
			"user1": {"admin"},
			"user2": {"viewer", "editor"},
			"root":  {"root"},
		},
		map[string][]string{
			"admin":  {"GET:/api/v1/users", "GET:/api/v1/users/:id"},
			"viewer": {"GET:/api/v1/posts"},
			"editor": {"POST:/api/v1/posts"},
		},
		[]string{
			"GET:/api/v1/users", "GET:/api/v1/users/:id",
			"GET:/api/v1/posts", "POST:/api/v1/posts",
			// 已注册但未授予任何角色 —— 覆盖"接口存在但无权访问"的拒绝路径
			"POST:/api/v1/users", "DELETE:/api/v1/posts",
		},
	)

	cases := []struct {
		name    string
		user    string
		path    string
		method  string
		wantErr bool
	}{
		{"allowed", "user1", "/api/v1/users", "GET", false},
		{"denied", "user1", "/api/v1/users", "POST", true},
		{"root bypass", "root", "/api/v1/anything", "DELETE", false},
		{"param path", "user1", "/api/v1/users/123", "GET", false},
		{"multiple roles allowed", "user2", "/api/v1/posts", "POST", false},
		{"multiple roles denied", "user2", "/api/v1/posts", "DELETE", true},
		{"empty user", "", "/api/v1/users", "GET", true},
		{"empty path", "user1", "", "GET", true},
		{"empty method", "user1", "/api/v1/users", "", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := e.Enforce(c.user, c.path, c.method)
			if (err != nil) != c.wantErr {
				t.Errorf("Enforce(%q, %q, %q) error=%v, wantErr=%v", c.user, c.path, c.method, err, c.wantErr)
			}
		})
	}
}
